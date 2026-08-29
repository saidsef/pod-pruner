package resources

import (
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func quietLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(io.Discard)
	return l
}

func TestForEachBoundedCapsConcurrency(t *testing.T) {
	const (
		items = 200
		limit = 10
	)

	var inFlight, visits int64
	var mu sync.Mutex
	peak := int64(0)

	work := make([]ContainerInfo, items)
	forEachBounded(work, limit, func(ContainerInfo) {
		atomic.AddInt64(&visits, 1)
		now := atomic.AddInt64(&inFlight, 1)
		mu.Lock()
		if now > peak {
			peak = now
		}
		mu.Unlock()
		time.Sleep(time.Millisecond)
		atomic.AddInt64(&inFlight, -1)
	})

	mu.Lock()
	got := peak
	mu.Unlock()

	if got > limit {
		t.Errorf("peak concurrency = %d, want at most %d", got, limit)
	}
	if got < 2 {
		t.Errorf("peak concurrency = %d, the work never ran in parallel", got)
	}
	if visits != items {
		t.Errorf("visited %d items, want %d", visits, items)
	}
	if left := atomic.LoadInt64(&inFlight); left != 0 {
		t.Errorf("%d still in flight after the call returned", left)
	}
	t.Logf("%d items, peak concurrency %d, cap %d", items, got, limit)
}

func TestForEachBoundedHandlesEmptyInput(t *testing.T) {
	called := false
	forEachBounded(nil, 10, func(ContainerInfo) { called = true })
	if called {
		t.Error("fn ran for an empty slice")
	}
}

func TestDeleteJobsDeletesEveryJob(t *testing.T) {
	client := fake.NewSimpleClientset()

	jobs := make([]ContainerInfo, 0, 25)
	for i := 0; i < 25; i++ {
		name := jobName(i)
		if _, err := client.BatchV1().Jobs("argocd").Create(t.Context(),
			&batchv1.Job{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "argocd"}}, metav1.CreateOptions{}); err != nil {
			t.Fatalf("seeding job: %v", err)
		}
		jobs = append(jobs, ContainerInfo{Namespace: "argocd", PodName: name, Status: "Complete"})
	}

	DeleteJobs(client, jobs, quietLogger())

	remaining, err := client.BatchV1().Jobs("argocd").List(t.Context(), metav1.ListOptions{})
	if err != nil {
		t.Fatalf("listing jobs: %v", err)
	}
	if len(remaining.Items) != 0 {
		t.Errorf("%d jobs left after DeleteJobs, want 0", len(remaining.Items))
	}
}

func jobName(i int) string {
	return "cron-job-" + string(rune('a'+i%26)) + string(rune('a'+(i/26)%26))
}
