package rsax

import (
	"os"
)

// LoadKeyHex loads an RSA key from a PEM file and returns the private and public keys as hex strings.
//
// Parameters:
//   - filename: The path to the PEM-encoded private key file.
//
// Returns:
//   - string: The hex-encoded private key.
//   - string: The hex-encoded public key.
//   - error: An error if reading or decoding fails.
func LoadKeyHex(filename string) (string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", "", err
	}
	return DecodeKeyHex(data)
}

// LoadKeyBase64 loads an RSA key from a PEM file and returns the private and public keys as Base64 strings.
//
// Parameters:
//   - filename: The path to the PEM-encoded private key file.
//
// Returns:
//   - string: The Base64-encoded private key.
//   - string: The Base64-encoded public key.
//   - error: An error if reading or decoding fails.
func LoadKeyBase64(filename string) (string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", "", err
	}
	return DecodeKeyBase64(data)
}
