package rsax

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/soyacen/gox/cryptox/hmacx"
	"golang.org/x/crypto/pkcs12"
	"net/url"
	"testing"
)

//go:embed data/test.p12
var privateKeyFromFile []byte

//go:embed data/test.cer
var publicKeyFromFile []byte

func TestRsa(t *testing.T) {
	// 原始数据
	info := "appClientId=test&detailStyle=1&id=com.miui.fm&nonce=-8162956393368422160:27234970&senderPackageName=com.android.browser"

	// 加签
	originalSign := hmacx.HmacSha256Hex([]byte("wcEvOZqgNtOURyNMdOm8Rg=="), []byte(info))
	fmt.Println("生成的签名:", originalSign)

	// 获取私钥
	privateKey, err := LoadPrivateKeyFromFile(privateKeyFromFile, "123456")
	if err != nil {
		panic(err)
	}

	// 获取公钥
	publicKey, err := LoadPublicKeyFromFile(publicKeyFromFile)
	if err != nil {
		panic(err)
	}

	// RSA加密
	encryptedResult, err := EncryptByPrivateKey([]byte(originalSign), publicKey)
	if err != nil {
		panic(err)
	}

	// Base64编码
	strBase64 := base64.StdEncoding.EncodeToString(encryptedResult)
	// URL编码
	sign := url.QueryEscape(strBase64)

	// 解码URL编码
	enResult, err := url.QueryUnescape(sign)
	if err != nil {
		panic(err)
	}

	// Base64解码
	decodedBytes, err := base64.StdEncoding.DecodeString(enResult)
	if err != nil {
		panic(err)
	}

	// 使用公钥解密
	decryptedSign, err := DecryptByPublicKey(decodedBytes, privateKey)
	if err != nil {
		panic(err)
	}

	// 判断解密结果与原始签名是否一致
	decodedSignStr := string(decryptedSign)
	fmt.Println("解密后的签名:", decodedSignStr)
	fmt.Println(decodedSignStr == originalSign)

}

const (
	KEY_SIZE     = 2048
	GROUP_SIZE   = KEY_SIZE / 8
	ENCRYPT_SIZE = GROUP_SIZE - 11
)

// EncryptByPrivateKey 使用私钥进行RSA加密
func EncryptByPrivateKey(data []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	cipherText := make([]byte, 0)
	for len(data) > 0 {
		blockSize := len(data)
		if blockSize > ENCRYPT_SIZE {
			blockSize = ENCRYPT_SIZE
		}
		ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data[:blockSize])
		if err != nil {
			return nil, err
		}
		cipherText = append(cipherText, ciphertext...)
		data = data[blockSize:]
	}
	return cipherText, nil
}

// DecryptByPublicKey 使用公钥进行RSA解密
func DecryptByPublicKey(encrypted []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	var plainText []byte
	for len(encrypted) > 0 {
		blockSize := len(encrypted)
		if blockSize > GROUP_SIZE {
			blockSize = GROUP_SIZE
		}
		plainTextBlock, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encrypted[:blockSize])
		if err != nil {
			return nil, err
		}
		encrypted = encrypted[blockSize:]
		plainText = append(plainText, plainTextBlock...)
	}
	return plainText, nil
}

// LoadPrivateKeyFromFile 加载PEM格式的私钥文件
func LoadPrivateKeyFromFile(pfxData []byte, password string) (*rsa.PrivateKey, error) {
	block, err := pkcs12.ToPEM(pfxData, password)
	if err != nil {
		return nil, err
	}
	if len(block) == 0 || block[0].Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing the private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block[0].Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// LoadPublicKeyFromFile 加载PEM格式的公钥证书文件
func LoadPublicKeyFromFile(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing the certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubkey := cert.PublicKey.(*rsa.PublicKey)
	return pubkey, nil
}
