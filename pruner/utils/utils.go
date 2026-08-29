/*
Copyright 2024 Said Sef

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

// GetEnv retrieves the value of the environment variable specified by key.
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

// GetEnvBool retrieves the environment variable specified by key and parses it
// as a boolean. If the variable is not set, or holds a value strconv.ParseBool
// cannot read, it returns defaultValue and logs the reason.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - defaultValue: The value to return if the variable is unset or unparseable.
// - log: A logger instance for logging warnings and errors.
//
// Returns:
// - The parsed value of the environment variable, or defaultValue.
func GetEnvBool(key string, defaultValue bool, log *logrus.Logger) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warnf("%s environment variable not set, defaulting to %t", key, defaultValue)
		return defaultValue
	}

	parsed, err := strconv.ParseBool(strings.TrimSpace(value))
	if err != nil {
		log.Errorf("%s environment variable is %q, which is not a boolean, defaulting to %t", key, value, defaultValue)
		return defaultValue
	}
	return parsed
}

// GetEnvDuration retrieves the environment variable specified by key and parses
// it as a Go duration, such as "90s" or "2m". If the variable is not set, holds
// a value time.ParseDuration cannot read, or is not positive, it returns
// defaultValue and logs the reason.
//
// Parameters:
// - key: The name of the environment variable to retrieve.
// - defaultValue: The value to return if the variable is unset or unusable.
// - log: A logger instance for logging warnings and errors.
//
// Returns:
// - The parsed duration, or defaultValue.
func GetEnvDuration(key string, defaultValue time.Duration, log *logrus.Logger) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Warnf("%s environment variable not set, defaulting to %s", key, defaultValue)
		return defaultValue
	}

	parsed, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		log.Errorf("%s environment variable is %q, which is not a duration, defaulting to %s", key, value, defaultValue)
		return defaultValue
	}
	if parsed <= 0 {
		log.Errorf("%s environment variable is %q, which is not positive, defaulting to %s", key, value, defaultValue)
		return defaultValue
	}
	return parsed
}

// SplitAndTrim splits a comma-separated value, trims surrounding whitespace
// from every entry and drops the entries that are left empty. An input that
// holds nothing but separators and whitespace yields an empty slice, so the
// caller can test for absence with len() rather than inspecting the elements.
//
// Parameters:
// - value: The raw comma-separated value to split.
//
// Returns:
// - The non-empty, whitespace-trimmed entries.
func SplitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	trimmed := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			trimmed = append(trimmed, part)
		}
	}
	return trimmed
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

// LogWithFields is a utility function for logging messages with different log levels.
// It logs the provided message along with any additional fields and an error if present.
//
// Parameters:
// - level: The log level at which to log the message (e.g., Error, Warn, Info, Debug).
// - fields: A map of fields to include in the log entry.
// - message: The message to log.
// - err: An optional error to include in the log entry.
//
// Returns:
// - None. The function logs the message at the specified log level.
func LogWithFields(level logrus.Level, fields []string, message string, errs ...error) {
	logFields := logrus.Fields{}

	// Convert []string to logrus.Fields
	for _, field := range fields {
		parts := strings.SplitN(field, ":", 2) // Split into key and value
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			logFields[key] = value
		}
	}

	LogWithMap(level, logFields, message, errs...)
}

// LogWithMap logs a message with fields supplied as a logrus.Fields map. Use it
// where a value is not a plain string, such as a list of resources or a count.
//
// Parameters:
// - level: The log level at which to log the message (e.g., Error, Warn, Info, Debug).
// - fields: The fields to attach to the log entry.
// - message: The message to log.
// - errs: Optional errors to include in the log entry.
//
// Returns:
// - None. The function logs the message at the specified log level.
func LogWithMap(level logrus.Level, fields logrus.Fields, message string, errs ...error) {
	if fields == nil {
		fields = logrus.Fields{}
	}

	// An error carries no exported fields, so the JSON formatter renders it as
	// {} and the message is lost. Stringify before adding it to the entry.
	if len(errs) > 0 {
		messages := make([]string, 0, len(errs))
		for _, err := range errs {
			if err != nil {
				messages = append(messages, err.Error())
			}
		}
		if len(messages) > 0 {
			fields["error"] = strings.Join(messages, "; ")
		}
	}

	entry := Logger().WithFields(fields)

	// Log based on the level
	switch level {
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	default:
		entry.Info(message)
	}
}

// Logger initializes and returns a singleton logrus Logger with JSON formatting.
// It ensures that only one instance of the logger is created using sync.Once.
// The logger is configured to use JSON formatting with timestamps enabled.
//
// Returns:
// *logrus.Logger: A singleton instance of the logrus Logger.
func Logger() *logrus.Logger {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: false})
	})
	return logger
}
