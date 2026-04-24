package auth

import (
	"encoding/base64"
	"strings"
)

// DecodeBasic decodes a Basic Authentication header string.
// It extracts the username and password from the authorization header.
//
// Parameters:
//   - auth: The authorization header string
//
// Returns:
//   - string: The decoded username
//   - string: The decoded password
//   - bool: True if decoding was successful
func DecodeBasic(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// EncodeBasic returns a Basic Authentication header string for the given username and password.
//
// Parameters:
//   - username: The username for authentication
//   - passwd: The password for authentication
//
// Returns:
//   - string: The Basic Authentication header value
func EncodeBasic(username, passwd string) string {
	auth := username + ":" + passwd
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
