package web

import (
	_ "embed"
	"fmt"
	"github.com/TMaize/dev-server/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed html/dir.html
var DirHtmlTemplate string

func RenderFile(w http.ResponseWriter, r *http.Request, file string) {
	if strings.HasSuffix(file, ".html") {
		RenderHtml(w, r, file)
	} else {
		http.ServeFile(w, r, file)
	}
}

func RenderDir(w http.ResponseWriter, r *http.Request, dir string) {
	indexFile := filepath.Join(dir, "index.html")
	_, err := os.Stat(indexFile)

	if err == nil {
		RenderHtml(w, r, indexFile)
		return
	}

	children, err := ioutil.ReadDir(dir)
	if err != nil {
		Render500(w, r, err)
		return
	}

	html := DirHtmlTemplate
	html = strings.ReplaceAll(html, `<title>Index Of /</title>`, fmt.Sprintf(`<title>Index Of %s</title>`, r.URL.Path))
	html = strings.ReplaceAll(html, `<h1>~/</h1>`, fmt.Sprintf(`<h1>~%s</h1>`, r.URL.Path))

	itemReg := regexp.MustCompile(`<a[\s\S]+?</a>`)
	items := itemReg.FindAllString(DirHtmlTemplate, -1)

	html = strings.ReplaceAll(html, items[0], "")
	html = strings.ReplaceAll(html, items[1], "")

	list := ""

	for _, child := range children {
		if child.IsDir() {
			item := items[0]
			item = strings.ReplaceAll(item, "dir name", child.Name())
			item = strings.ReplaceAll(item, `href=""`, fmt.Sprintf(`href="%s/"`, url.PathEscape(child.Name())))
			list += item + "\n"
		}
	}

	for _, child := range children {
		if !child.IsDir() {
			item := items[1]
			item = strings.ReplaceAll(item, "file name", child.Name())
			item = strings.ReplaceAll(item, "file size", util.FmtFileSize(child.Size()))
			item = strings.ReplaceAll(item, `href=""`, fmt.Sprintf(`href="%s"`, url.PathEscape(child.Name())))
			list += item + "\n"
		}
	}

	html = strings.ReplaceAll(html, `<div class="list">`, `<div class="list">`+list)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(html))
}

func RenderHtml(w http.ResponseWriter, r *http.Request, file string) {
	http.ServeFile(w, r, file)
}

func Render404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func Render500(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, fmt.Sprintf("Server Error: %s", err), http.StatusInternalServerError)
}
