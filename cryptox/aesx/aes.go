package aesx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Cipher defines the interface for AES encryption and decryption operations.
type Cipher interface {
	// Encrypt encrypts plaintext using the provided key.
	//
	// Parameters:
	//   - plaintext: the data to encrypt
	//   - key: the encryption key
	//
	// Returns:
	//   - []byte: the encrypted ciphertext
	//   - error: any error encountered during encryption
	Encrypt(plaintext, key []byte) ([]byte, error)

	// Decrypt decrypts ciphertext using the provided key.
	//
	// Parameters:
	//   - ciphertext: the data to decrypt
	//   - key: the decryption key
	//
	// Returns:
	//   - []byte: the decrypted plaintext
	//   - error: any error encountered during decryption
	Decrypt(ciphertext, key []byte) ([]byte, error)
}

// ECB returns a Cipher that uses Electronic Codebook (ECB) mode.
//
// Returns:
//   - Cipher: an ECB mode cipher instance
func ECB() Cipher {
	return ecb{}
}

// CBC returns a Cipher that uses Cipher Block Chaining (CBC) mode.
//
// Returns:
//   - Cipher: a CBC mode cipher instance
func CBC() Cipher {
	return cbc{}
}

// CFB returns a Cipher that uses Cipher Feedback (CFB) mode.
//
// Returns:
//   - Cipher: a CFB mode cipher instance
func CFB() Cipher {
	return cfb{}
}

// OFB returns a Cipher that uses Output Feedback (OFB) mode.
//
// Returns:
//   - Cipher: an OFB mode cipher instance
func OFB() Cipher {
	return ofb{}
}

// CTR returns a Cipher that uses Counter (CTR) mode.
//
// Returns:
//   - Cipher: a CTR mode cipher instance
func CTR() Cipher {
	return ctr{}
}

type ecb struct {
}

// Encrypt encrypts plaintext using ECB mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (ecb) Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 填充明文
	plaintext = padding(plaintext, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	// 分块加密
	for i := 0; i < len(plaintext); i += block.BlockSize() {
		block.Encrypt(ciphertext[i:i+block.BlockSize()], plaintext[i:i+block.BlockSize()])
	}
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using ECB mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (ecb) Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plaintext := make([]byte, len(ciphertext))
	// 分块解密
	for i := 0; i < len(ciphertext); i += block.BlockSize() {
		block.Decrypt(plaintext[i:i+block.BlockSize()], ciphertext[i:i+block.BlockSize()])
	}
	// 去除填充
	plaintext, err = unPadding(plaintext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

type cbc struct{}

// Encrypt encrypts plaintext using CBC mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (cbc) Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plaintext = padding(plaintext, blockSize)
	ciphertext := make([]byte, blockSize+len(plaintext))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using CBC mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (cbc) Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext, err = unPadding(plaintext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

type cfb struct{}

// Encrypt encrypts plaintext using CFB mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (cfb) Encrypt(plaintext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewCFBEncrypter}.Encrypt(plaintext, key)
}

// Decrypt decrypts ciphertext using CFB mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (cfb) Decrypt(ciphertext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewCFBDecrypter}.Decrypt(ciphertext, key)
}

type ofb struct {
}

// Encrypt encrypts plaintext using OFB mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (ofb) Encrypt(plaintext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewOFB}.Encrypt(plaintext, key)
}

// Decrypt decrypts ciphertext using OFB mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (ofb) Decrypt(ciphertext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewOFB}.Decrypt(ciphertext, key)
}

type ctr struct{}

// Encrypt encrypts plaintext using CTR mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (ctr) Encrypt(plaintext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewCTR}.Encrypt(plaintext, key)
}

// Decrypt decrypts ciphertext using CTR mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (ctr) Decrypt(ciphertext, key []byte) ([]byte, error) {
	return baseStream{newStream: cipher.NewCTR}.Decrypt(ciphertext, key)
}

type baseStream struct {
	newStream func(block cipher.Block, iv []byte) cipher.Stream
}

// Encrypt encrypts plaintext using a stream cipher mode.
//
// Parameters:
//   - plaintext: the data to encrypt
//   - key: the encryption key
//
// Returns:
//   - []byte: the encrypted ciphertext
//   - error: any error encountered during encryption
func (b baseStream) Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := b.newStream(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using a stream cipher mode.
//
// Parameters:
//   - ciphertext: the data to decrypt
//   - key: the decryption key
//
// Returns:
//   - []byte: the decrypted plaintext
//   - error: any error encountered during decryption
func (b baseStream) Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	plaintext := make([]byte, len(ciphertext))
	stream := b.newStream(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext, nil
}

// padding 填充数据
func padding(data []byte, blockSize int) []byte {
	size := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(size)}, size)...)
}

// unPadding 去除填充
func unPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	size := int(data[length-1])
	if size > length {
		return nil, fmt.Errorf("invalid size")
	}
	return data[:length-size], nil
}
