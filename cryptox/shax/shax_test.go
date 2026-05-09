package shax

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestSha1(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha1.Sum(tc.data)
			actual := Sha1(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha1(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha1Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha1Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha224(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha256.Sum224(tc.data)
			actual := Sha224(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha224(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha224Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha224Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha256.Sum256(tc.data)
			actual := Sha256(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha256(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha256Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha256Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha384(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha512.Sum384(tc.data)
			actual := Sha384(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha384(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha384Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha384Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha512(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha512.Sum512(tc.data)
			actual := Sha512(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha512(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha512Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha512Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha512_224(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha512.Sum512_224(tc.data)
			actual := Sha512_224(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha512_224(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha512_224Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha512_224Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}

func TestSha512_256(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{name: "empty", data: []byte{}},
		{name: "hello world", data: []byte("hello world")},
		{name: "special chars", data: []byte("!@#$%^&*()")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expected := sha512.Sum512_256(tc.data)
			actual := Sha512_256(tc.data)
			if !bytes.Equal(actual, expected[:]) {
				t.Errorf("Sha512_256(%q) = %x, want %x", tc.data, actual, expected)
			}

			expectedHex := hex.EncodeToString(expected[:])
			actualHex := Sha512_256Hex(tc.data)
			if actualHex != expectedHex {
				t.Errorf("Sha512_256Hex(%q) = %s, want %s", tc.data, actualHex, expectedHex)
			}
		})
	}
}
