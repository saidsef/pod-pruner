package resources

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetContainers retrieves a list of container names in a given namespace that are in specified states.
// It checks the environment variable CONTAINER_STATUSES to determine which states to filter by.
//
// Parameters:
// - clientset: A pointer to the Kubernetes clientset.
// - namespace: The namespace from which to retrieve the containers.
//
// Returns:
// - A slice of strings containing the names of the containers in the specified states.
// - An error if the retrieval of pods fails.
func GetContainers(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	statuses := strings.Split(os.Getenv("CONTAINER_STATUSES"), ",")
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var containers []string
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if IsContainerInState(containerStatus, statuses) {
				containers = append(containers, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
			}
		}
	}
	return containers, nil
}

// IsContainerInState checks if a container's status matches any of the specified states.
//
// Parameters:
// - containerStatus: The status of the container to check.
// - statuses: A slice of strings representing the states to check against.
//
// Returns:
// - A boolean indicating whether the container is in one of the specified states.
func IsContainerInState(containerStatus v1.ContainerStatus, statuses []string) bool {
	if containerStatus.State.Waiting != nil && utils.Contains(statuses, containerStatus.State.Waiting.Reason) {
		return true
	}
	if containerStatus.State.Terminated != nil && utils.Contains(statuses, containerStatus.State.Terminated.Reason) {
		return true
	}
	return false
}

// DeleteContainers deletes the specified containers in a given namespace.
// It logs the success or failure of each deletion attempt.
//
// Parameters:
// - clientset: A pointer to the Kubernetes clientset.
// - namespace: The namespace from which to delete the containers.
// - containers: A slice of strings representing the containers to delete.
// - log: A logger instance for logging the results of the deletion attempts.
func DeleteContainers(clientset *kubernetes.Clientset, namespace string, containers []string, log *logrus.Logger) {
	for _, container := range containers {
		parts := strings.Split(container, ": ")
		if len(parts) != 2 {
			log.Warnf("Unexpected format for container state: '%s'", container)
			continue
		}
		podInfo := parts[0]
		containerName := strings.Split(podInfo, "/")[1]

		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), containerName, metav1.DeleteOptions{})
		if err != nil {
			log.Errorf("Failed to delete pod '%s' in namespace '%s': %v", podInfo, namespace, err)
		} else {
			log.Infof("Successfully deleted pod '%s' in namespace '%s'", podInfo, namespace)
		}
	}
}
