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

package resources

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetContainers retrieves a list of container names from pods in the specified namespace
// that are in the states defined by the CONTAINER_STATUSES environment variable.
// It returns a slice of container names in the format "namespace/podName: containerName".
// If the environment variable is not set or is empty, an error is returned.
// If there is an error while listing the pods, it returns an error with context.
//
// Parameters:
// - clientset: A Kubernetes clientset used to interact with the Kubernetes API.
// - namespace: The namespace from which to retrieve the pods.
//
// Returns:
// - A slice of strings containing the names of the containers in the specified states.
// - An error if the environment variable is not set, empty, or if there is an error
// while listing the pods.
func GetContainers(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	statuses := strings.Split(os.Getenv("CONTAINER_STATUSES"), ",")
	if len(statuses) == 0 || (len(statuses) == 1 && statuses[0] == "") {
		return nil, fmt.Errorf("CONTAINER_STATUSES environment variable is not set or empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var containers []string
	var continueToken string

	for {
		podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list pods in namespace '%s': %w", namespace, err)
		}

		for _, pod := range podList.Items {
			for _, containerStatus := range pod.Status.ContainerStatuses {
				if isContainerInState(containerStatus, statuses) {
					containers = append(containers, fmt.Sprintf("%s/%s: %s", pod.Namespace, pod.Name, containerStatus.Name))
				}
			}
		}

		if podList.Continue == "" {
			break
		}
		continueToken = podList.Continue
	}

	return containers, nil
}

// isContainerInState checks if the given container status is in one of the specified states.
// It returns true if the container is waiting or terminated with a reason that matches one of the statuses.
//
// Parameters:
// - containerStatus: The status of the container to check.
// - statuses: A slice of strings representing the states to check against.
//
// Returns:
// - A boolean indicating whether the container status matches one of the specified states.
func isContainerInState(containerStatus v1.ContainerStatus, statuses []string) bool {
	statusSet := make(map[string]struct{}, len(statuses))
	for _, status := range statuses {
		statusSet[status] = struct{}{}
	}

	if containerStatus.State.Waiting != nil {
		if _, exists := statusSet[containerStatus.State.Waiting.Reason]; exists {
			return true
		}
	}
	if containerStatus.State.Terminated != nil {
		if _, exists := statusSet[containerStatus.State.Terminated.Reason]; exists {
			return true
		}
	}
	return false
}

// DeleteContainers deletes the specified containers (pods) in the given namespace.
// It logs warnings for any containers that do not conform to the expected format.
// If a pod deletion fails, it logs an error; otherwise, it logs a success message.
//
// Parameters:
// - clientset: A Kubernetes clientset used to interact with the Kubernetes API.
// - namespace: The namespace from which to delete the pods.
// - containers: A slice of strings containing the names of the containers to delete.
// - log: A logger used to log messages regarding the deletion process.
func DeleteContainers(clientset *kubernetes.Clientset, namespace string, containers []string, log *logrus.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, container := range containers {
		parts := strings.Split(container, ": ")
		if len(parts) != 2 {
			log.WithFields(logrus.Fields{
				"state": container,
				"parts": parts,
			}).Warn("Unexpected format for container")
			continue
		}
		podInfo := parts[0]
		podParts := strings.Split(podInfo, "/")
		if len(podParts) != 2 {
			log.WithFields(logrus.Fields{
				"podInfo": podInfo,
			}).Warn("Unexpected format for pod info")
			continue
		}
		podName := podParts[1]

		err := clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{})
		if err != nil {
			log.WithFields(logrus.Fields{
				"pod":       podInfo,
				"namespace": namespace,
				"error":     err,
			}).Error("Failed to delete pod")
		} else {
			log.WithFields(logrus.Fields{
				"pod":       podInfo,
				"namespace": namespace,
			}).Info("Successfully deleted pod")
		}
	}
}
