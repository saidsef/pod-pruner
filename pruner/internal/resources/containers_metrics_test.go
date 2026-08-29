package resources

import (
	"io"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/saidsef/pod-pruner/pruner/internal/metrics"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// counterValue reads a counter without pulling in the prometheus test helper,
// which needs a module the build does not otherwise use.
func counterValue(t *testing.T, c prometheus.Counter) float64 {
	t.Helper()
	var m dto.Metric
	if err := c.Write(&m); err != nil {
		t.Fatalf("reading counter: %v", err)
	}
	return m.GetCounter().GetValue()
}

func discardLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestDeleteContainersCountsPrunedPods(t *testing.T) {
	client := fake.NewSimpleClientset()

	targets := []ContainerInfo{
		{Namespace: "argocd", PodName: "pod-a", Status: "Completed"},
		{Namespace: "argocd", PodName: "pod-b", Status: "Completed"},
		{Namespace: "argocd", PodName: "pod-c", Status: "Error"},
	}
	for _, target := range targets {
		if _, err := client.CoreV1().Pods(target.Namespace).Create(t.Context(),
			&v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: target.PodName, Namespace: target.Namespace}},
			metav1.CreateOptions{}); err != nil {
			t.Fatalf("seeding pod: %v", err)
		}
	}

	before := counterValue(t, metrics.PodsPruned.WithLabelValues("argocd", "Completed"))
	DeleteContainers(client, targets, discardLogger())

	if got := counterValue(t, metrics.PodsPruned.WithLabelValues("argocd", "Completed")) - before; got != 2 {
		t.Errorf("pods_pruned_total{state=\"Completed\"} rose by %v, want 2", got)
	}
	if got := counterValue(t, metrics.PodsPruned.WithLabelValues("argocd", "Error")); got != 1 {
		t.Errorf("pods_pruned_total{state=\"Error\"} = %v, want 1", got)
	}

	remaining, err := client.CoreV1().Pods("argocd").List(t.Context(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("listing pods: %v", err)
	}
	if len(remaining.Items) != 0 {
		t.Errorf("%d pods left after DeleteContainers, want 0", len(remaining.Items))
	}
}

func TestDeleteContainersLeavesMetricAloneOnFailure(t *testing.T) {
	client := fake.NewSimpleClientset()

	before := counterValue(t, metrics.PodsPruned.WithLabelValues("argocd", "Unknown"))
	DeleteContainers(client, []ContainerInfo{{Namespace: "argocd", PodName: "absent", Status: "Unknown"}}, discardLogger())

	if got := counterValue(t, metrics.PodsPruned.WithLabelValues("argocd", "Unknown")) - before; got != 0 {
		t.Errorf("counter rose by %v on a failed delete, want 0", got)
	}
}
