package tlsx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

// NewServerTLSConfig creates a new TLS configuration for a server.
//
// Parameters:
//   - certPath: path to the server certificate file, or empty string to use a random certificate
//   - keyPath: path to the server key file, or empty string to use a random certificate
//   - caPath: path to the CA certificate file, or empty string to skip client certificate verification
//
// Returns:
//   - *tls.Config: the created TLS configuration
//   - error: any error encountered during certificate loading or generation
func NewServerTLSConfig(certPath, keyPath, caPath string) (*tls.Config, error) {
	res := &tls.Config{}
	if certPath == "" || keyPath == "" {
		cert, err := NewRandomCert()
		if err != nil {
			return nil, err
		}
		res.Certificates = []tls.Certificate{cert}
	} else {
		cert, err := NewCustomCert(certPath, keyPath)
		if err != nil {
			return nil, err
		}
		res.Certificates = []tls.Certificate{cert}
	}
	if caPath == "" {
		return res, nil
	}
	pool, err := NewCertPool(caPath)
	if err != nil {
		return nil, err
	}
	res.ClientAuth = tls.RequireAndVerifyClientCert
	res.ClientCAs = pool
	return res, nil
}

// NewClientTLSConfig creates a new TLS configuration for a client.
//
// Parameters:
//   - certPath: path to the client certificate file, or empty string to skip client certificate
//   - keyPath: path to the client key file, or empty string to skip client certificate
//   - caPath: path to the CA certificate file, or empty string to skip server certificate verification
//   - serverName: the expected server name for SNI verification
//
// Returns:
//   - *tls.Config: the created TLS configuration
//   - error: any error encountered during certificate loading
func NewClientTLSConfig(certPath, keyPath, caPath, serverName string) (*tls.Config, error) {
	res := &tls.Config{}
	if certPath != "" && keyPath != "" {
		cert, err := NewCustomCert(certPath, keyPath)
		if err != nil {
			return nil, err
		}
		res.Certificates = []tls.Certificate{cert}
	}
	res.ServerName = serverName
	if caPath != "" {
		pool, err := NewCertPool(caPath)
		if err != nil {
			return nil, err
		}

		res.RootCAs = pool
		res.InsecureSkipVerify = false
	} else {
		res.InsecureSkipVerify = true
	}
	return res, nil
}

// NewRandomCert generates a random TLS certificate and returns it as a tls.Certificate.
//
// Returns:
//   - tls.Certificate: the generated random certificate
//   - error: any error encountered during key generation or certificate creation
func NewRandomCert() (tls.Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&key.PublicKey,
		key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return tlsCert, nil
}

// NewCustomCert loads a TLS certificate from the specified paths and returns it as a tls.Certificate.
//
// Parameters:
//   - certPath: path to the certificate file
//   - keyPath: path to the key file
//
// Returns:
//   - tls.Certificate: the loaded certificate
//   - error: any error encountered during certificate loading
func NewCustomCert(certPath, keyPath string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certPath, keyPath)
}

// NewCertPool creates a new x509.CertPool and adds the CA certificate from the specified path.
//
// Parameters:
//   - caPath: path to the CA certificate file
//
// Returns:
//   - *x509.CertPool: the created certificate pool containing the CA certificate
//   - error: any error encountered while reading the CA certificate file
func NewCertPool(caPath string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	caCrt, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)
	return pool, nil
}

// NewRandomPrivateKey generates a random RSA private key and returns it as a PEM-encoded byte slice.
//
// Returns:
//   - []byte: the PEM-encoded private key
//   - error: any error encountered during key generation
func NewRandomPrivateKey() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}), nil
}
