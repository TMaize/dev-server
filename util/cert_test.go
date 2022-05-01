package util

import (
	"io/ioutil"
	"net"
	"os"
	"testing"
)

func TestGenerateCACert(t *testing.T) {
	cer, key, err := GenerateCACert()
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("D:/ca.cer", cer, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("D:/ca.key", key, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateCertByCA(t *testing.T) {

	caKeyData, err := ioutil.ReadFile("D:/ca.key")
	if err != nil {
		t.Fatal(err)
	}

	caCerData, err := ioutil.ReadFile("D:/ca.cer")
	if err != nil {
		t.Fatal(err)
	}

	alternateIPs := []net.IP{net.IPv4(127, 0, 0, 1), net.IPv4(127, 0, 0, 2)}
	alternateDNS := []string{"dev-server.land"}
	cert, key, err := GenerateCertByCA(caKeyData, caCerData, "localhost", alternateIPs, alternateDNS)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("D:/out.cer", cert, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile("D:/out.key", key, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

