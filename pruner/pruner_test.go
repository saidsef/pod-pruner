package main

import (
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func quietLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestSweepPrunesJobsWhenPodListingFails(t *testing.T) {
	// An empty CONTAINER_STATUSES makes the pod half of the sweep fail.
	t.Setenv("CONTAINER_STATUSES", "")
	t.Setenv("JOB_STATUSES", "Complete")

	client := fake.NewSimpleClientset(&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: "finished", Namespace: "batch"},
		Status: batchv1.JobStatus{Conditions: []batchv1.JobCondition{
			{Type: batchv1.JobComplete, Status: v1.ConditionTrue},
		}},
	})

	sweep(client, []string{"batch"}, []string{"PODS", "JOBS"}, false, quietLogger())

	jobs, err := client.BatchV1().Jobs("batch").List(t.Context(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("listing jobs: %v", err)
	}
	if len(jobs.Items) != 0 {
		t.Errorf("%d jobs left after the sweep, want 0", len(jobs.Items))
	}
}
