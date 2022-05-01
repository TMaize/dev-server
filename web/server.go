package web

import (
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"github.com/TMaize/dev-server/util"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type StaticServer struct {
	Https       bool
	Address     string
	Port        uint
	Domain      string
	Root        string
	caFile      string
	cerData     []byte
	keyData     []byte
	certificate tls.Certificate
	server      *http.Server
}

func (s *StaticServer) PreRun() error {
	if s.Port == 0 {
		s.Port = 80
		if s.Https {
			s.Port = 443
		}
	}

	if s.Https && s.Port == 80 {
		return errors.New("can't listen 80 for https")
	}

	if s.Root == "" {
		s.Root = "."
	}
	s.Root = util.FmtFilePath(s.Root)

	// init ca cer
	if s.Https {
		if err := util.InstallCACer(); err != nil {
			return err
		}
		cerFile, _ := util.GetConfigFile("ca.cer")
		s.caFile = cerFile
	}

	// init domain cer
	if s.Https {
		alternateIPs := []net.IP{net.IPv4(127, 0, 0, 1)}
		alternateDNS := make([]string, 0)
		if s.Domain != "localhost" {
			alternateDNS = append(alternateDNS, s.Domain)
		}
		if s.Address != "127.0.0.1" {
			alternateIPs = append(alternateIPs, net.ParseIP(s.Address).To4())
		}
		cerByte, keyByte, err := util.GenerateCertByDefaultCA("localhost", alternateIPs, alternateDNS)
		if err != nil {
			return errors.New("GenerateCertByDefaultCA Error: " + err.Error())
		}

		pair, err := tls.X509KeyPair(cerByte, keyByte)
		if err != nil {
			return errors.New("X509KeyPair Error: " + err.Error())
		}

		s.cerData = cerByte
		s.keyData = keyByte
		s.certificate = pair
	}

	return nil
}

func (s *StaticServer) PrintArgs() {

	urlList := make([]string, 0)

	urlList = append(urlList, util.BuildURL(s.Https, s.Address, s.Port))
	urlList = append(urlList, util.BuildURL(s.Https, s.Domain, s.Port))

	fmt.Printf("  https: %v\n", s.Https)
	fmt.Printf("address: %s\n", s.Address)
	fmt.Printf("   port: %d\n", s.Port)
	fmt.Printf("   root: %s\n", s.Root)
	if s.Https {
		fmt.Printf(" use CA: %s\n", s.caFile)
	}
	fmt.Printf("    url: %v\n", strings.Join(urlList, " , "))

	time.Sleep(time.Second * 3)
	fmt.Println("press ctrl-c to stop.")
}

func (s *StaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file := path.Join(s.Root, r.URL.Path)
	info, err := os.Stat(file)

	if r.Method == "GET" {
		w.Header().Set("Cache-Control", "no-cache")
	}

	if err != nil {
		if os.IsNotExist(err) {
			Render404(w, r)
		} else {
			Render500(w, r, err)
		}
		return
	}

	if strings.HasSuffix(r.URL.Path, "/") && info.IsDir() {
		RenderDir(w, r, file)
		return
	}

	if !strings.HasSuffix(r.URL.Path, "/") && !info.IsDir() {
		RenderFile(w, r, file)
		return
	}

	Render404(w, r)
}

func (s *StaticServer) Run() error {
	if err := s.PreRun(); err != nil {
		return err
	}

	// custom init
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler: s,
	}

	go s.PrintArgs()

	if s.Https {
		s.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{s.certificate},
		}
		return s.server.ListenAndServeTLS("", "")
	}

	return s.server.ListenAndServe()
}
