package rsax

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
)

// DecodeKeyHex decodes a PEM-encoded RSA private key and returns the private and public keys as hex strings.
//
// Parameters:
//   - data: The PEM-encoded RSA private key data.
//
// Returns:
//   - string: The hex-encoded private key.
//   - string: The hex-encoded public key.
//   - error: An error if decoding fails.
func DecodeKeyHex(data []byte) (string, string, error) {
	block, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", "", err
	}
	privateKeyStr := hex.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	publicKeyStr := hex.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))
	return privateKeyStr, publicKeyStr, nil
}

// DecodeKeyBase64 decodes a PEM-encoded RSA private key and returns the private and public keys as Base64 strings.
//
// Parameters:
//   - data: The PEM-encoded RSA private key data.
//
// Returns:
//   - string: The Base64-encoded private key.
//   - string: The Base64-encoded public key.
//   - error: An error if decoding fails.
func DecodeKeyBase64(data []byte) (string, string, error) {
	block, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", "", err
	}
	privateKeyStr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	publicKeyStr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))
	return privateKeyStr, publicKeyStr, nil
}
