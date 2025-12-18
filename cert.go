package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

func generateCert() error {
	psk, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			Organization: []string{"nosni"},
			CommonName: "nosni",
		},
		NotBefore: time.Now(),
		NotAfter: time.Now().AddDate(8, 0, 0),
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &psk.PublicKey, psk)
	if err != nil {
		return err
	}

	certFile, err := os.Create("ca.crt")
	if err != nil {
		return err
	}
	defer certFile.Close()

	err = pem.Encode(certFile, &pem.Block{
		Type: "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return err
	}

	keyFile, err := os.Create("key.pem")
	if err != nil {
		return err
	}
	defer keyFile.Close()

	err = pem.Encode(keyFile, &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(psk),
	})
	if err != nil {
		return err
	}

	return nil
}

func loadCert() (*tls.Certificate, error) {
	_, err := os.Stat("key.pem")
	if err != nil {
		if os.IsNotExist(err) {
			err := generateCert()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	cert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	key, err := os.ReadFile("key.pem")
	if err != nil {
		return nil, err
	}

	parsedCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	parsedCert.Leaf, err = x509.ParseCertificate(parsedCert.Certificate[0])
	if err != nil {
		return nil, err
	}

	return &parsedCert, nil
}