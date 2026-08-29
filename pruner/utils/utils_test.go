package utils

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"io"
)

func quietLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestGetEnvBool(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		set      bool
		fallback bool
		want     bool
	}{
		{"unset falls back", "", false, true, true},
		{"lower true", "true", true, true, true},
		{"lower false", "false", true, true, false},
		{"title case True", "True", true, true, true},
		{"upper TRUE", "TRUE", true, true, true},
		{"title case False", "False", true, true, false},
		{"numeric 1", "1", true, true, true},
		{"numeric 0", "0", true, true, false},
		{"padded false", " false ", true, true, false},
		{"padded true", " true ", true, true, true},
		{"typo falls back to dry run", "flase", true, true, true},
		{"yes is not a bool", "yes", true, true, true},
		{"empty falls back", "", true, true, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Unsetenv("DRY_RUN")
			if c.set {
				os.Setenv("DRY_RUN", c.value)
				defer os.Unsetenv("DRY_RUN")
			}
			if got := GetEnvBool("DRY_RUN", c.fallback, quietLogger()); got != c.want {
				t.Errorf("GetEnvBool(%q) = %t, want %t", c.value, got, c.want)
			}
		})
	}
}
