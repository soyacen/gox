package auth

import "testing"

func TestEncodeDecodeBasic(t *testing.T) {
	const (
		username = "testuser"
		password = "testpass"
	)

	encoded := EncodeBasic(username, password)
	decodedUser, decodedPass, ok := DecodeBasic(encoded)

	if !ok {
		t.Fatal("failed to decode basic auth")
	}
	if decodedUser != username {
		t.Fatalf("username mismatch: got %s, want %s", decodedUser, username)
	}
	if decodedPass != password {
		t.Fatalf("password mismatch: got %s, want %s", decodedPass, password)
	}
}

func TestDecodeBasicInvalid(t *testing.T) {
	testCases := []string{
		"",
		"NotBasic",
		"Basic invalid-base64",
		"Basic dXNlcjEyMw==", // Base64 of "user123", no colon
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			_, _, ok := DecodeBasic(tc)
			if ok {
				t.Fatal("expected decode to fail")
			}
		})
	}
}
