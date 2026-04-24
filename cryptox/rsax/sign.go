package rsax

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"

	"github.com/soyacen/gox/cryptox/shax"
)

// SignWithSha256Hex signs data using an RSA private key with SHA-256 and returns the signature as a hex string.
// The private key must be provided as a hex-encoded PKCS#1 RSA private key.
//
// Parameters:
//   - data: The data to sign.
//   - priKey: The hex-encoded RSA private key.
//
// Returns:
//   - string: The hex-encoded signature.
//   - error: An error if signing fails.
func SignWithSha256Hex(data []byte, priKey string) (string, error) {
	priBytes, err := hex.DecodeString(priKey)
	if err != nil {
		return "", err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(priBytes)
	if err != nil {
		return "", err
	}
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, shax.Sha256(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sign), nil
}

// VerifySignWithSha256Hex verifies an RSA SHA-256 signature using a hex-encoded public key.
// The signature and public key must be provided as hex-encoded strings.
//
// Parameters:
//   - data: The original data that was signed.
//   - hexSign: The hex-encoded signature.
//   - hexPubKey: The hex-encoded RSA public key.
//
// Returns:
//   - error: An error if verification fails.
func VerifySignWithSha256Hex(data []byte, hexSign, hexPubKey string) error {
	sig, err := hex.DecodeString(hexSign)
	if err != nil {
		return err
	}
	pubBytes, err := hex.DecodeString(hexPubKey)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKCS1PublicKey(pubBytes)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, shax.Sha256(data), sig)
}

// SignWithSha256Base64 signs data using an RSA private key with SHA-256 and returns the signature as a Base64 string.
// The private key must be provided as a Base64-encoded PKCS#1 RSA private key.
//
// Parameters:
//   - data: The data to sign.
//   - priKey: The Base64-encoded RSA private key.
//
// Returns:
//   - string: The Base64-encoded signature.
//   - error: An error if signing fails.
func SignWithSha256Base64(data []byte, priKey string) (string, error) {
	der, err := base64.StdEncoding.DecodeString(priKey)
	if err != nil {
		return "", err
	}
	priv, err := x509.ParsePKCS1PrivateKey(der)
	if err != nil {
		return "", err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, shax.Sha256(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// VerifySignWithSha256Base64 verifies an RSA SHA-256 signature using a Base64-encoded public key.
// The signature and public key must be provided as Base64-encoded strings.
//
// Parameters:
//   - data: The original data that was signed.
//   - base64Sign: The Base64-encoded signature.
//   - base64PubKey: The Base64-encoded RSA public key.
//
// Returns:
//   - error: An error if verification fails.
func VerifySignWithSha256Base64(data []byte, base64Sign, base64PubKey string) error {
	sig, err := base64.StdEncoding.DecodeString(base64Sign)
	if err != nil {
		return err
	}
	der, err := base64.StdEncoding.DecodeString(base64PubKey)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKCS1PublicKey(der)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, shax.Sha256(data), sig)
}
