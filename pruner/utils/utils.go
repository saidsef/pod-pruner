package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// getEnv retrieves the value of the environment variable specified by key.
// If the variable is not set, it returns the defaultValue and logs a warning.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - defaultValue: The value to return if the environment variable is not set.
// - log: A logger instance for logging warnings.
//
// Returns:
// - The value of the environment variable or the default value if not set.
func GetEnv(key, defaultValue string, log *logrus.Logger) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warnf("%s environment variable not set, defaulting to %s", key, defaultValue)
		return defaultValue
	}
	return value
}

// Contains checks if a string is present in a slice of strings.
//
// Parameters:
// - list: A slice of strings to search through.
// - str: The string to search for.
//
// Returns:
// - A boolean indicating whether the string is found in the slice.
func Contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
