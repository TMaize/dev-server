package util

import (
	"bytes"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

// LoadPrivateKey rsa.PrivateKey
// -----BEGIN RSA PRIVATE KEY-----
// ...
// -----END RSA PRIVATE KEY-----
func LoadPrivateKey(keyData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyData)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// LoadCer x509.Certificate
// -----BEGIN CERTIFICATE-----
// ...
// -----END CERTIFICATE-----
func LoadCer(cerData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(cerData)
	cer, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cer, nil
}

// GenerateCACert CA cer&key
// 根证书带域名也是可以的，但是最好还是分开
// CA证书信任后，后面所有签发的证书都是被信任的
func GenerateCACert() ([]byte, []byte, error) {
	validFrom := time.Now().Add(-time.Hour)
	maxAge := time.Hour * 24 * 365 * 10

	caKey, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "CA-DEV",
			Organization: []string{"DEV-SERVER"},
			Country:      []string{"US"},
		},
		NotBefore: validFrom,
		NotAfter:  validFrom.Add(maxAge),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDERBytes, err := x509.CreateCertificate(cryptorand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate cert
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: caDERBytes}); err != nil {
		return nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(caKey)}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

// GenerateCertByCA cer&key
func GenerateCertByCA(caKeyData []byte, caCerData []byte, host string, alternateIPs []net.IP, alternateDNS []string) ([]byte, []byte, error) {
	caKey, err := LoadPrivateKey(caKeyData)
	if err != nil {
		return nil, nil, err
	}
	caCer, err := LoadCer(caCerData)
	if err != nil {
		return nil, nil, err
	}

	validFrom := time.Now().Add(-time.Hour)
	maxAge := time.Hour * 24 * 30

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			CommonName: host,
		},
		NotBefore: validFrom,
		NotAfter:  validFrom.Add(maxAge),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	priv, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, caCer, &priv.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	// Generate cert, followed by ca
	certBuffer := bytes.Buffer{}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}
	if err := pem.Encode(&certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: caCer.Raw}); err != nil {
		return nil, nil, err
	}

	// Generate key
	keyBuffer := bytes.Buffer{}
	if err := pem.Encode(&keyBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}

func GenerateCertByDefaultCA(host string, alternateIPs []net.IP, alternateDNS []string) ([]byte, []byte, error) {
	cerFile, err := GetConfigFile("ca.cer")
	if err != nil {
		return nil, nil, err
	}

	keyFile, err := GetConfigFile("ca.key")
	if err != nil {
		return nil, nil, err
	}

	caCerData, err := ioutil.ReadFile(cerFile)
	if err != nil {
		return nil, nil, err
	}

	caKeyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	return GenerateCertByCA(caKeyData, caCerData, host, alternateIPs, alternateDNS)
}

func InstallCACer() error {
	cerFile, err := GetConfigFile("ca.cer")
	if err != nil {
		return err
	}

	keyFile, err := GetConfigFile("ca.key")
	if err != nil {
		return err
	}

	_, err = os.Stat(cerFile)

	// install once
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	cer, key, err := GenerateCACert()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cerFile, cer, os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(keyFile, key, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
