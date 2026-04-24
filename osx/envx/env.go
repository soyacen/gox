package envx

import (
	"os"
	"strconv"
)

// GetString retrieves the value of the specified environment variable.
// If the environment variable is not set, it returns the provided default value.
//
// Parameters:
//   - key: The name of the environment variable to look up.
//   - defaultValue: The fallback value to return if the environment variable is not set.
//
// Returns:
//   - string: The value of the environment variable if set, otherwise the defaultValue.
func GetString(key, defaultValue string) string {
	// Attempt to lookup the environment variable
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defaultValue
}

// GetInt retrieves an integer value from environment variables with a fallback default value.
//
// Parameters:
//   - key: The environment variable key to look up.
//   - defaultValue: The fallback value to return if the key is not found or conversion fails.
//
// Returns:
//   - int: The retrieved integer value or defaultValue if not found or invalid.
//   - error: Conversion error if environment variable exists but is not a valid integer.
func GetInt(key string, defaultValue int) (int, error) {
	// Attempt to lookup environment variable
	v, ok := os.LookupEnv(key)
	if ok {
		// Convert string to integer if found
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}

	// Return default value if environment variable doesn't exist
	return defaultValue, nil
}

// GetFloat64 retrieves a float64 value from environment variables.
// If the environment variable exists and can be parsed as float64, returns the parsed value.
// If the environment variable does not exist, returns the provided default value.
// If the environment variable exists but cannot be parsed as float64, returns default value and an error.
//
// Parameters:
//   - key: The name of the environment variable to look up.
//   - defaultValue: The default value to return if the environment variable is not set.
//
// Returns:
//   - float64: The retrieved value or default value.
//   - error: Parsing error if the environment variable exists but cannot be parsed as float64.
func GetFloat64(key string, defaultValue float64) (float64, error) {
	// Check if environment variable exists
	v, ok := os.LookupEnv(key)
	if ok {
		// Attempt to parse the environment variable value as float64
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// GetBool retrieves a boolean value from environment variables with a fallback default.
// It first checks if the specified environment variable exists. If it exists, it attempts
// to parse the value as a boolean. If parsing fails, it returns the error along with the
// default value. If the environment variable doesn't exist, it returns the default value.
//
// Parameters:
//   - key: The name of the environment variable to look up.
//   - defaultValue: The value to return if the environment variable is not set.
//
// Returns:
//   - bool: The parsed boolean value from environment variable or the default value.
//   - error: Parsing error if environment variable exists but cannot be parsed as boolean.
func GetBool(key string, defaultValue bool) (bool, error) {
	// Check if environment variable exists
	v, ok := os.LookupEnv(key)
	if ok {
		// Parse environment variable value as boolean
		value, err := strconv.ParseBool(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}
