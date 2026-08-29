package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

// capture redirects the singleton logger and returns the decoded entry.
func capture(t *testing.T, emit func()) map[string]any {
	t.Helper()
	var buf bytes.Buffer
	Logger().SetOutput(&buf)
	defer Logger().SetOutput(nil)
	emit()

	entry := map[string]any{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v (%q)", err, buf.String())
	}
	return entry
}

func TestLogWithMapKeepsListsAndCounts(t *testing.T) {
	names := []string{"argocd/pod-a (Completed)", "argocd/pod-b (Error)"}
	entry := capture(t, func() {
		LogWithMap(logrus.InfoLevel, logrus.Fields{"count": len(names), "containers": names},
			"Dry run mode. The following containers would be deleted")
	})

	if entry["count"] != float64(2) {
		t.Errorf("count = %v, want 2", entry["count"])
	}
	got, ok := entry["containers"].([]any)
	if !ok {
		t.Fatalf("containers = %#v, want a list", entry["containers"])
	}
	if len(got) != 2 || got[0] != names[0] || got[1] != names[1] {
		t.Errorf("containers = %v, want %v", got, names)
	}
}

func TestLogWithFieldsParsesKeyValuePairs(t *testing.T) {
	entry := capture(t, func() {
		LogWithFields(logrus.InfoLevel, []string{"pod:my-pod", "namespace:argocd"}, "deleted")
	})

	if entry["pod"] != "my-pod" || entry["namespace"] != "argocd" {
		t.Errorf("fields = %v, want pod=my-pod namespace=argocd", entry)
	}
}

func TestLogWithMapStringifiesErrors(t *testing.T) {
	entry := capture(t, func() {
		LogWithMap(logrus.ErrorLevel, logrus.Fields{"pod": "my-pod"}, "failed", errBoom, nil, errBang)
	})

	if entry["error"] != "boom; bang" {
		t.Errorf("error = %v, want \"boom; bang\"", entry["error"])
	}
}

var (
	errBoom = errors.New("boom")
	errBang = errors.New("bang")
)
