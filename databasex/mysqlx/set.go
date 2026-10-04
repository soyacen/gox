package mysqlx

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// Set represents a MySQL SET value, a list of zero or more members chosen from a
// predefined list. MySQL renders a SET as a comma separated string and members
// cannot contain a comma.
type Set struct {
	// Strings are the selected members. A non-NULL value with no selected
	// members is represented by an empty, non-nil slice.
	Strings []string
	// Valid reports whether the value is not NULL.
	Valid bool
}

// Scan implements the database/sql.Scanner interface.
//
// Parameters:
//   - src: a value produced by the database driver.
//
// Returns:
//   - error: non-nil when src cannot be scanned into a Set.
func (dst *Set) Scan(src any) error {
	v, isNull, err := coerceString(src)
	if err != nil {
		return err
	}

	if isNull {
		*dst = Set{}
		return nil
	}

	*dst = Set{Strings: splitSet(v), Valid: true}
	return nil
}

// Value implements the database/sql/driver.Valuer interface. The members are
// joined with commas.
//
// Returns:
//   - driver.Value: the comma separated member list, or nil when the value is
//     NULL.
//   - error: always nil.
func (src Set) Value() (driver.Value, error) {
	if !src.Valid {
		return nil, nil
	}
	return strings.Join(src.Strings, ","), nil
}

// MarshalJSON implements the encoding/json.Marshaler interface. The value is
// encoded as a JSON array of strings.
//
// Returns:
//   - []byte: the JSON array, or null when the value is NULL.
//   - error: non-nil when the value cannot be marshaled.
func (src Set) MarshalJSON() ([]byte, error) {
	if !src.Valid {
		return []byte("null"), nil
	}
	if src.Strings == nil {
		return []byte("[]"), nil
	}
	return json.Marshal(src.Strings)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface. The value
// is expected to be a JSON array of strings.
//
// Parameters:
//   - data: the raw JSON document.
//
// Returns:
//   - error: non-nil when data is not a JSON array of strings or null.
func (dst *Set) UnmarshalJSON(data []byte) error {
	if isJSONNull(data) {
		*dst = Set{}
		return nil
	}

	var members []string
	if err := json.Unmarshal(data, &members); err != nil {
		return err
	}
	if members == nil {
		members = []string{}
	}

	*dst = Set{Strings: members, Valid: true}
	return nil
}

// splitSet splits a MySQL SET literal into its members. The empty literal
// yields an empty, non-nil slice.
//
// Parameters:
//   - s: the comma separated SET literal.
//
// Returns:
//   - []string: the selected members.
func splitSet(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
