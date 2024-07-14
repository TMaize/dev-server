package web

import (
	"crypto/tls"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/TMaize/dev-server/util"
)

type Server struct {
	Https       bool
	Address     string
	Port        uint
	Domain      string
	Root        string
	caFile      string
	cerData     []byte
	keyData     []byte
	certificate tls.Certificate
	ipList      []string
}

func (s *Server) PreRun() error {
	if s.Port == 0 {
		s.Port = 80
		if s.Https {
			s.Port = 443
		}
	}

	if s.Address == "0.0.0.0" {
		list := util.GetLocalIP()
		s.ipList = append(s.ipList, list...)
	} else {
		s.ipList = append(s.ipList, s.Address)
	}

	if !s.Https && s.Port == 443 {
		return errors.New("can't listen 443 for http")
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
		alternateIPs := make([]net.IP, 0) //[]net.IP{net.IPv4(127, 0, 0, 1)}
		alternateDNS := make([]string, 0)

		if s.Domain != "localhost" {
			alternateDNS = append(alternateDNS, s.Domain)
		}

		for _, addr := range s.ipList {
			alternateIPs = append(alternateIPs, net.ParseIP(addr).To4())
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

func (s *Server) PrintArgs() {

	urlList := make([]string, 0)

	for _, addr := range s.ipList {
		urlList = append(urlList, util.BuildURL(s.Https, addr, s.Port))
	}
	urlList = append(urlList, util.BuildURL(s.Https, s.Domain, s.Port))

	fmt.Printf("  https: %v\n", s.Https)
	fmt.Printf("address: %s\n", s.Address)
	fmt.Printf("   port: %d\n", s.Port)
	fmt.Printf("   root: %s\n", s.Root)
	if s.Https {
		fmt.Printf(" use CA: %s\n", s.caFile)
	}
	fmt.Printf("    url: %v\n", urlList[0])

	for _, item := range urlList[1:] {
		fmt.Printf("         %v\n", item)
	}

	time.Sleep(time.Second * 3)
	fmt.Println("press ctrl-c to stop.")
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) Run() error {
	if err := s.PreRun(); err != nil {
		return err
	}

	// custom init
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler: s,
	}

	go s.PrintArgs()

	if s.Https {
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{s.certificate},
		}
		return server.ListenAndServeTLS("", "")
	}

	// ipv4 and ipv6
	//return server.ListenAndServe()

	// ipv4
	ln, err := net.Listen("tcp4", server.Addr)
	if err != nil {
		return err
	}
	return server.Serve(ln)
}
