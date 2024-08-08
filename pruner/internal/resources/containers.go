package resources

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetContainers retrieves a list of container names from pods in the specified namespace
// that are in the states defined by the CONTAINER_STATUSES environment variable.
// It returns a slice of container names in the format "namespace/podName: containerName".
// If the environment variable is not set or empty, an error is returned.
// If there is an error while listing the pods, it returns an error with context.
func GetContainers(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	statuses := strings.Split(os.Getenv("CONTAINER_STATUSES"), ",")
	if len(statuses) == 0 {
		return nil, fmt.Errorf("CONTAINER_STATUSES environment variable is not set or empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
	}

	var containers []string
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if isContainerInState(containerStatus, statuses) {
				containers = append(containers, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
			}
		}
	}
	return containers, nil
}

// isContainerInState checks if the given container status is in one of the specified states.
// It returns true if the container is waiting or terminated with a reason that matches one of the statuses.
func isContainerInState(containerStatus v1.ContainerStatus, statuses []string) bool {
	if containerStatus.State.Waiting != nil && utils.Contains(statuses, containerStatus.State.Waiting.Reason) {
		return true
	}
	if containerStatus.State.Terminated != nil && utils.Contains(statuses, containerStatus.State.Terminated.Reason) {
		return true
	}
	return false
}

// DeleteContainers deletes the specified containers (pods) in the given namespace.
// It logs warnings for any containers that do not conform to the expected format.
// If a pod deletion fails, it logs an error; otherwise, it logs a success message.
func DeleteContainers(clientset *kubernetes.Clientset, namespace string, containers []string, log *logrus.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, container := range containers {
		parts := strings.Split(container, ": ")
		if len(parts) != 2 {
			log.Warnf("Unexpected format for container state: '%s'", container)
			continue
		}
		podInfo := parts[0]
		podParts := strings.Split(podInfo, "/")
		if len(podParts) != 2 {
			log.Warnf("Unexpected format for pod info: '%s'", podInfo)
			continue
		}
		podName := podParts[1]

		err := clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
		if err != nil {
			log.Errorf("Failed to delete pod '%s' in namespace '%s': %v", podInfo, namespace, err)
		} else {
			log.Infof("Successfully deleted pod '%s' in namespace '%s'", podInfo, namespace)
		}
	}
}
