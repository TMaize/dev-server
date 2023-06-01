package web

import (
	"encoding/json"
	"github.com/TMaize/dev-server/ui"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

type RespBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RSF interface {
	io.Reader
	io.Seeker
	io.Closer
	Stat() (fs.FileInfo, error)
}

func RenderContent(w http.ResponseWriter, r *http.Request, files []string) {
	for _, file := range files {
		fileLoc := file
		var err error
		var f fs.File

		if strings.HasPrefix(file, "embed:") {
			fileLoc = file[6:]
			f, err = ui.FS.Open(ui.Prefix + fileLoc)
		} else {
			f, err = os.Open(fileLoc)
		}

		if err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// TODO warning
		defer func(f fs.File) {
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}(f)

		stat, err := f.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			continue
		}

		http.ServeContent(w, r, fileLoc, stat.ModTime(), f.(io.ReadSeeker))
		return
	}

	http.Error(w, "404 page not found", http.StatusNotFound)
}

func RenderJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	body, _ := json.Marshal(data)
	_, _ = w.Write(body)
}

func RenderInfo(s *Server, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/..") {
		RenderJSON(w, RespBody{Code: 400, Message: "invalid URL path"}, 200)
		return
	}

	file := path.Join(s.Root, r.URL.Path)
	info, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			RenderJSON(w, RespBody{Code: 404, Message: "no such file or directory"}, 200)
			return
		}
		if strings.Contains(err.Error(), "not a directory") {
			RenderJSON(w, RespBody{Code: 404, Message: "no such file or directory"}, 200)
			return
		}
		RenderJSON(w, RespBody{Code: 500, Message: err.Error()}, 200)
		return
	}

	if !info.IsDir() {
		RenderJSON(w, RespBody{Code: 200, Data: File{Name: info.Name(), Size: info.Size(), Type: "file"}}, 200)
		return
	}
	children, err := os.ReadDir(file)
	if err != nil {
		RenderJSON(w, RespBody{Code: 500, Message: err.Error()}, 200)
		return
	}

	files := make([]File, 0)

	for _, child := range children {
		if child.IsDir() {
			files = append(files, File{Name: child.Name(), Size: 0, Type: "dir"})
		} else {
			fileInfo, _ := child.Info()
			files = append(files, File{Name: child.Name(), Size: fileInfo.Size(), Type: "file"})
		}
	}

	RenderJSON(w, RespBody{Code: 200, Data: files}, 200)
}

func RenderStatic(s *Server, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/..") {
		http.Error(w, "invalid URL path", 400)
		return
	}

	// 重写
	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path = r.URL.Path + "index.html"
		RenderStatic(s, w, r)
		return
	}

	file := path.Join(s.Root, r.URL.Path)
	info, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			if s.UI {
				RenderContent(w, r, []string{"embed:" + r.URL.Path, "embed:/index.html"})
			} else {
				http.Error(w, "404 page not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 目录必须以 xxx/ 的格式访问
	if info.IsDir() {
		if s.UI {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		} else {
			http.Error(w, "404 page not found", http.StatusNotFound)
		}
		return
	}

	// 禁止在 ServeFile 重定向 /xxx/index.html -> /xxx/
	if strings.HasSuffix(r.URL.Path, "/index.html") {
		if s.UI {
			RenderContent(w, r, []string{file, "embed:" + r.URL.Path})
		} else {
			RenderContent(w, r, []string{file})
		}
		return
	}

	http.ServeFile(w, r, file)
}
