package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/TMaize/dev-server/util"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Server struct {
	Https       bool
	Address     string
	Port        uint
	Domain      string
	Cors        bool
	Target      string
	caFile      string
	cerData     []byte
	keyData     []byte
	certificate tls.Certificate
}

func (s *Server) director(req *http.Request) {
	target, _ := url.Parse(s.Target)

	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.Host = target.Host

	req.Header["X-Forwarded-For"] = nil
}

func (s *Server) modify(resp *http.Response) error {
	if s.Cors {
		resp.Header.Set("Access-Control-Allow-Origin", resp.Request.Header.Get("Origin"))
		resp.Header.Set("Access-Control-Allow-Methods", "*")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Expose-Headers", "Token")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		resp.Header.Set("Access-Control-Max-Age", "3600")
	}

	return nil
}

func (s *Server) PreRun() error {
	if s.Port == 0 {
		s.Port = 80
		if s.Https {
			s.Port = 443
		}
	}

	if !s.Https && s.Port == 443 {
		return errors.New("can't listen 443 for http")
	}

	if s.Https && s.Port == 80 {
		return errors.New("can't listen 80 for https")
	}

	targetUrl, err := url.Parse(s.Target)
	if err != nil || (targetUrl.Scheme != "http" && targetUrl.Scheme != "https") {
		return errors.New("Invalid target: " + s.Target)
	}

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

func (s *Server) PrintArgs() {
	urlList := make([]string, 0)

	urlList = append(urlList, util.BuildURL(s.Https, s.Address, s.Port))
	urlList = append(urlList, util.BuildURL(s.Https, s.Domain, s.Port))

	fmt.Printf("  https: %v\n", s.Https)
	fmt.Printf("address: %s\n", s.Address)
	fmt.Printf("   port: %d\n", s.Port)
	if s.Https {
		fmt.Printf(" use CA: %s\n", s.caFile)
	}
	fmt.Printf("  proxy: [%s] => %s\n", strings.Join(urlList, ", "), s.Target)

	time.Sleep(time.Second * 3)
	fmt.Println("press ctrl-c to stop.")
}

func (s *Server) Run() error {
	if err := s.PreRun(); err != nil {
		return err
	}

	handler := httputil.ReverseProxy{
		Director:       s.director,
		ModifyResponse: s.modify,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler: &handler,
	}

	go s.PrintArgs()

	if s.Https {
		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{s.certificate},
		}
		return server.ListenAndServeTLS("", "")
	}

	return server.ListenAndServe()
}
