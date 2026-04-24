package rsax

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
)

// EncryptToHex encrypts plaintext using an RSA public key and returns the ciphertext as a hex string.
// The public key must be provided as a hex-encoded PKCS#1 RSA public key.
//
// Parameters:
//   - plainText: The data to encrypt.
//   - hexPubKey: The hex-encoded RSA public key.
//
// Returns:
//   - string: The hex-encoded ciphertext.
//   - error: An error if decryption fails.
func EncryptToHex(plainText []byte, hexPubKey string) (string, error) {
	pub, err := hex.DecodeString(hexPubKey)
	if err != nil {
		return "", err
	}
	cipherBytes, err := rsaEncrypt(plainText, pub)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipherBytes), nil
}

// DecryptByHex decrypts a hex-encoded ciphertext using an RSA private key.
// The private key must be provided as a hex-encoded PKCS#1 RSA private key.
//
// Parameters:
//   - hexCipherText: The hex-encoded ciphertext.
//   - hexPriKey: The hex-encoded RSA private key.
//
// Returns:
//   - []byte: The decrypted plaintext.
//   - error: An error if decryption fails.
func DecryptByHex(hexCipherText, hexPriKey string) ([]byte, error) {
	privateBytes, err := hex.DecodeString(hexPriKey)
	if err != nil {
		return nil, err
	}
	cipherTextBytes, err := hex.DecodeString(hexCipherText)
	if err != nil {
		return nil, err
	}
	return rsaDecrypt(cipherTextBytes, privateBytes)
}

// EncryptToBase64 encrypts plaintext using an RSA public key and returns the ciphertext as a Base64 string.
// The public key must be provided as a Base64-encoded PKCS#1 RSA public key.
//
// Parameters:
//   - plainText: The data to encrypt.
//   - base64PubKey: The Base64-encoded RSA public key.
//
// Returns:
//   - string: The Base64-encoded ciphertext.
//   - error: An error if encryption fails.
func EncryptToBase64(plainText []byte, base64PubKey string) (string, error) {
	pub, err := base64.StdEncoding.DecodeString(base64PubKey)
	if err != nil {
		return "", err
	}
	cipherBytes, err := rsaEncrypt(plainText, pub)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// DecryptByBase64 decrypts a Base64-encoded ciphertext using an RSA private key.
// The private key must be provided as a Base64-encoded PKCS#1 RSA private key.
//
// Parameters:
//   - base64CipherText: The Base64-encoded ciphertext.
//   - base64PriKey: The Base64-encoded RSA private key.
//
// Returns:
//   - []byte: The decrypted plaintext.
//   - error: An error if decryption fails.
func DecryptByBase64(base64CipherText, base64PriKey string) ([]byte, error) {
	privateBytes, err := base64.StdEncoding.DecodeString(base64PriKey)
	if err != nil {
		return nil, err
	}
	cipherTextBytes, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return nil, err
	}
	return rsaDecrypt(cipherTextBytes, privateBytes)
}

func rsaEncrypt(plainText, publicKey []byte) ([]byte, error) {
	pub, err := x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pubSize, plainTextSize := pub.Size(), len(plainText)
	offSet, once := 0, pubSize-11
	buffer := bytes.Buffer{}
	for offSet < plainTextSize {
		endIndex := offSet + once
		if endIndex > plainTextSize {
			endIndex = plainTextSize
		}
		bytesOnce, err := rsa.EncryptPKCS1v15(rand.Reader, pub, plainText[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}

func rsaDecrypt(cipherText, privateKey []byte) ([]byte, error) {
	pri, err := x509.ParsePKCS1PrivateKey(privateKey)
	if err != nil {
		return []byte{}, err
	}
	priSize, cipherTextSize := pri.Size(), len(cipherText)
	offSet := 0
	buffer := bytes.Buffer{}
	for offSet < cipherTextSize {
		endIndex := offSet + priSize
		if endIndex > cipherTextSize {
			endIndex = cipherTextSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, pri, cipherText[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}
