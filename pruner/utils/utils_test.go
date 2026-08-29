package utils

import (
	"os"
	"testing"
	"time"

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

func TestSplitAndTrim(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  []string
	}{
		{"empty is no entries", "", nil},
		{"whitespace only is no entries", "   ", nil},
		{"separators only is no entries", ",,,", nil},
		{"single entry", "argocd", []string{"argocd"}},
		{"padded entries", "argocd, tekton", []string{"argocd", "tekton"}},
		{"trailing separator", "Error,", []string{"Error"}},
		{"inner blank dropped", "Error,,Completed", []string{"Error", "Completed"}},
		{"tabs and newlines", "Error,\tCompleted\n", []string{"Error", "Completed"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := SplitAndTrim(c.value)
			if len(got) != len(c.want) {
				t.Fatalf("SplitAndTrim(%q) = %q, want %q", c.value, got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("SplitAndTrim(%q) = %q, want %q", c.value, got, c.want)
				}
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	const fallback = 120 * time.Second

	cases := []struct {
		name  string
		value string
		set   bool
		want  time.Duration
	}{
		{"unset falls back", "", false, fallback},
		{"seconds", "90s", true, 90 * time.Second},
		{"minutes", "2m", true, 2 * time.Minute},
		{"padded", " 45s ", true, 45 * time.Second},
		{"bare number is not a duration", "60", true, fallback},
		{"nonsense falls back", "soon", true, fallback},
		{"empty falls back", "", true, fallback},
		{"zero falls back", "0s", true, fallback},
		{"negative falls back", "-30s", true, fallback},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			os.Unsetenv("INTERVAL")
			if c.set {
				os.Setenv("INTERVAL", c.value)
				defer os.Unsetenv("INTERVAL")
			}
			if got := GetEnvDuration("INTERVAL", fallback, quietLogger()); got != c.want {
				t.Errorf("GetEnvDuration(%q) = %s, want %s", c.value, got, c.want)
			}
		})
	}
}

func TestNamespaceLabel(t *testing.T) {
	cases := []struct {
		name      string
		namespace string
		want      string
	}{
		{"named namespace", "argocd", "argocd"},
		{"empty is every namespace", "", "*"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NamespaceLabel(c.namespace); got != c.want {
				t.Errorf("NamespaceLabel(%q) = %q, want %q", c.namespace, got, c.want)
			}
		})
	}
}
