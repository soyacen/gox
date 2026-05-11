package rsax

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignHex(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyHex(1024)
	assert.NoError(t, err)

	t.Log(privateKey)
	t.Log(publicKey)

	rawData := []byte("he is hello kitty")
	signHex, err := SignWithSha256Hex(rawData, privateKey)
	assert.NoError(t, err)

	err = VerifySignWithSha256Hex(rawData, signHex, publicKey)
	assert.NoError(t, err)
}

func TestSignBase64(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyBase64(1024)
	assert.NoError(t, err)

	t.Log(privateKey)
	t.Log(publicKey)

	rawData := []byte("he is hello kitty")
	signHex, err := SignWithSha256Base64(rawData, privateKey)
	assert.NoError(t, err)

	err = VerifySignWithSha256Base64(rawData, signHex, publicKey)
	assert.NoError(t, err)
}

func TestCrypt(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyHex(1024)
	assert.NoError(t, err)

	rawData := []byte("he is hello kitty")
	encryptDate, err := EncryptToHex(rawData, publicKey)
	assert.NoError(t, err)

	data, err := DecryptByHex(encryptDate, privateKey)
	assert.NoError(t, err)

	assert.Equal(t, rawData, data)
}

func TestLoad(t *testing.T) {
	// 临时目录生成 PEM 私钥文件,避免依赖 /tmp/priv.pem
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	file := filepath.Join(t.TempDir(), "priv.pem")
	assert.NoError(t, os.WriteFile(file, pemBytes, 0o600))

	priv, pub, err := LoadKeyBase64(file)
	assert.NoError(t, err)
	assert.NotEmpty(t, priv)
	assert.NotEmpty(t, pub)
	t.Log(priv)
	t.Log(pub)

	privHex, pubHex, err := LoadKeyHex(file)
	assert.NoError(t, err)
	assert.NotEmpty(t, privHex)
	assert.NotEmpty(t, pubHex)
}

func TestLoad_MissingFile(t *testing.T) {
	_, _, err := LoadKeyBase64(filepath.Join(t.TempDir(), "does-not-exist.pem"))
	assert.Error(t, err)
}
