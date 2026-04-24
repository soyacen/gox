package cryptox

import (
	"crypto/x509"
	"encoding/pem"
)

// CertToPEMString takes an x509 certificate and returns a PEM-encoded string.
//
// Parameters:
//   - cert: the x509 certificate to encode
//
// Returns:
//   - string: the PEM-encoded certificate
func CertToPEMString(cert *x509.Certificate) string {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return string(pem.EncodeToMemory(&block))
}
