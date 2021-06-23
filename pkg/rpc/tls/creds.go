package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

// NewClientCreds impl
func NewClientCreds(svrName, caFile, certFile, keyFile string) (*credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(ca) {
		return nil, errors.New("rootCAs append failed")
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
		ServerName:   svrName,
	})
	return &creds, nil
}

// NewServerCreds impl
func NewServerCreds(caFile, certFile, keyFile string) (*credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	rootCAs := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	if !rootCAs.AppendCertsFromPEM(ca) {
		return nil, errors.New("rootCAs append failed")
	}
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    rootCAs,
	})
	return &creds, nil
}
