package rsax

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
)

// GenerateKeyHex generates a new RSA key pair and returns the private and public keys as hex strings.
//
// Parameters:
//   - bits: The size of the RSA key in bits (e.g., 2048).
//
// Returns:
//   - string: The hex-encoded private key.
//   - string: The hex-encoded public key.
//   - error: An error if key generation fails.
func GenerateKeyHex(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	privateKeyStr := hex.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	publicKeyStr := hex.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))
	return privateKeyStr, publicKeyStr, nil
}

// GenerateKeyBase64 generates a new RSA key pair and returns the private and public keys as Base64 strings.
//
// Parameters:
//   - bits: The size of the RSA key in bits (e.g., 2048).
//
// Returns:
//   - string: The Base64-encoded private key.
//   - string: The Base64-encoded public key.
//   - error: An error if key generation fails.
func GenerateKeyBase64(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	privateKeyStr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey))
	publicKeyStr := base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&privateKey.PublicKey))
	return privateKeyStr, publicKeyStr, nil
}
