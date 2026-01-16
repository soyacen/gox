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
