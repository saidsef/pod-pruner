package resources

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func terminated(name, reason string) v1.ContainerStatus {
	return v1.ContainerStatus{
		Name:  name,
		State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{Reason: reason}},
	}
}

func waiting(name, reason string) v1.ContainerStatus {
	return v1.ContainerStatus{
		Name:  name,
		State: v1.ContainerState{Waiting: &v1.ContainerStateWaiting{Reason: reason}},
	}
}

func running(name string) v1.ContainerStatus {
	return v1.ContainerStatus{Name: name, State: v1.ContainerState{Running: &v1.ContainerStateRunning{}}}
}

func pod(statuses ...v1.ContainerStatus) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "my-pod", Namespace: "argocd"},
		Status:     v1.PodStatus{ContainerStatuses: statuses},
	}
}

func TestPruneCandidateYieldsOneEntryPerPod(t *testing.T) {
	statuses := []string{"Completed", "Error"}

	cases := []struct {
		name       string
		pod        v1.Pod
		wantMatch  bool
		wantReason string
	}{
		{"single matching container", pod(terminated("app", "Completed")), true, "Completed"},
		{"sidecar pod matches once", pod(terminated("app", "Completed"), terminated("istio-proxy", "Completed")), true, "Completed"},
		{"three matching containers match once", pod(terminated("a", "Completed"), terminated("b", "Completed"), terminated("c", "Completed")), true, "Completed"},
		{"reason taken from the first match", pod(running("app"), terminated("sidecar", "Error")), true, "Error"},
		{"waiting reason matches", pod(waiting("app", "Error")), true, "Error"},
		{"no containers no match", pod(), false, ""},
		{"running pod does not match", pod(running("app")), false, ""},
		{"unlisted reason does not match", pod(terminated("app", "OOMKilled")), false, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, matched := pruneCandidate(c.pod, statuses)
			if matched != c.wantMatch {
				t.Fatalf("matched = %t, want %t", matched, c.wantMatch)
			}
			if !matched {
				return
			}
			if got.PodName != "my-pod" || got.Namespace != "argocd" {
				t.Errorf("got %s/%s, want argocd/my-pod", got.Namespace, got.PodName)
			}
			if got.Status != c.wantReason {
				t.Errorf("status = %q, want %q", got.Status, c.wantReason)
			}
		})
	}
}
