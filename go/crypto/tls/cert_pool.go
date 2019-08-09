package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

// LoadCertificatePool returns loads a TLS x509.CertPool or update a TLS x509.CertPool if nil.
// certString: Base64 encoded (without padding) string of the TLS certificate (PEM encoded) to be used for HTTP over TLS (HTTPS).
// Example: certString="-----BEGIN CERTIFICATE-----\nMIIDZTCCAk2gAwIBAgIEV5xOtDANBgkqhkiG9w0BAQ0FADA0MTIwMAYDVQQDDClP..."
// certPath: The path to the TLS certificate (pem encoded).
// Example: certPath=~/cert.pem
// certs: certs of x509.Certificate, tls.Certificate, *x509.Certificate, *tls.Certificate
func LoadCertificatePool(
	certPool *x509.CertPool,
	certString string,
	certFile string,
	certs ...interface{},
) (*x509.CertPool, error) {
	var tlsCertBytes []byte
	var err error
	if certString == "" && certFile == "" && len(certs) == 0 {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	} else if certString != "" {
		tlsCertBytes, err = base64.StdEncoding.DecodeString(certString)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
		}
	} else if certFile != "" {
		tlsCertBytes, err = ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
	} else {
		var loaded bool
		for _, cert := range certs {
			if certPool == nil {
				certPool = x509.NewCertPool()
			}
			switch cert.(type) {
			case *x509.Certificate:
				x509Cert := cert.(*x509.Certificate)
				certPool.AddCert(x509Cert)
				loaded = true
			case x509.Certificate:
				x509Cert := cert.(x509.Certificate)
				certPool.AddCert(&x509Cert)
				loaded = true
			case *tls.Certificate:
				tlsCert := cert.(*tls.Certificate)
				for _, certBytes := range tlsCert.Certificate {
					x509Cert, err := x509.ParseCertificate(certBytes)
					if err != nil {
						continue
					}
					certPool.AddCert(x509Cert)
					loaded = true
				}
			case tls.Certificate:
				tlsCert := cert.(tls.Certificate)
				for _, certBytes := range tlsCert.Certificate {
					x509Cert, err := x509.ParseCertificate(certBytes)
					if err != nil {
						continue
					}
					certPool.AddCert(x509Cert)
					loaded = true
				}
			}
		}
		if loaded {
			return certPool, nil
		}
	}

	if len(tlsCertBytes) == 0 {
		return nil, errors.WithStack(ErrInvalidCertificateConfiguration)
	}
	if certPool == nil {
		certPool = x509.NewCertPool()
	}
	if !certPool.AppendCertsFromPEM(tlsCertBytes) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return certPool, nil
}
