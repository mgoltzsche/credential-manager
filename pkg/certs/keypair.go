package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

const (
	pemTypeCertificate   = "CERTIFICATE"
	pemTypeRSAPrivateKey = "RSA PRIVATE KEY"
)

var (
	caTTL   = 24 * 365 * 5 * time.Hour
	certTTL = 24 * 90 * time.Hour
)

type KeyPair struct {
	caCertPEM     []byte
	signedCertPEM []byte
	x509Cert      *x509.Certificate
	rsaKey        *rsa.PrivateKey
}

func X509KeyPair(keyPEM, certPEM, caCertPEM []byte) (*KeyPair, error) {
	pb, _ := pem.Decode(keyPEM)
	rsaKey, err := x509.ParsePKCS1PrivateKey(pb.Bytes)
	if err != nil {
		return nil, err
	}
	pb, _ = pem.Decode(certPEM)
	cert, err := x509.ParseCertificate(pb.Bytes)
	if err != nil {
		return nil, err
	}
	return &KeyPair{rsaKey: rsaKey, x509Cert: cert, signedCertPEM: certPEM, caCertPEM: caCertPEM}, nil
}

func (p *KeyPair) CACertPEM() []byte {
	return p.caCertPEM
}

func (p *KeyPair) CertPEM() []byte {
	return p.signedCertPEM
}

func (p *KeyPair) KeyPEM() []byte {
	privKeyPEM := new(bytes.Buffer)
	pem.Encode(privKeyPEM, &pem.Block{
		Type:  pemTypeRSAPrivateKey,
		Bytes: x509.MarshalPKCS1PrivateKey(p.rsaKey),
	})
	return privKeyPEM.Bytes()
}

func (p *KeyPair) NeedsRenewal() bool {
	cert := p.x509Cert
	ttl := cert.NotAfter.Sub(cert.NotBefore)
	renewalTime := cert.NotAfter.Add(ttl / -4)
	return time.Now().After(renewalTime)
}

func NewSelfSignedCAKeyPair(commonName string) (*KeyPair, error) {
	now := time.Now()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			Organization: []string{"Fake org"},
			CommonName:   commonName,
		},
		NotBefore:             now,
		NotAfter:              now.Add(caTTL),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	return genKeyPair(ca, 4096, nil)
}

func NewServerKeyPair(dnsNames []string, ca *KeyPair) (*KeyPair, error) {
	if len(dnsNames) == 0 {
		return nil, errors.New("gen server key pair: no dns names provided")
	}
	now := time.Now()
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		DNSNames:     dnsNames,
		Subject: pkix.Name{
			CommonName: dnsNames[0],
		},
		NotBefore:             now,
		NotAfter:              now.Add(certTTL),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	return genKeyPair(cert, 2048, ca)
}

func genKeyPair(cert *x509.Certificate, bits int, ca *KeyPair) (*KeyPair, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	signKey := privKey
	caCert := cert
	var caCertPEM []byte
	if ca != nil {
		caCert = ca.x509Cert
		signKey = ca.rsaKey
		caCertPEM = ca.CertPEM()
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &privKey.PublicKey, signKey)
	if err != nil {
		return nil, err
	}
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  pemTypeCertificate,
		Bytes: certBytes,
	})
	if ca == nil {
		caCertPEM = certPEM.Bytes()
	}
	return &KeyPair{
		rsaKey:        privKey,
		x509Cert:      cert,
		signedCertPEM: certPEM.Bytes(),
		caCertPEM:     caCertPEM,
	}, nil
}
