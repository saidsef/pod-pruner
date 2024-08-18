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
	"strings"
	"sync"

	"github.com/saidsef/pod-pruner/pruner/internal/metrics"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetJobs retrieves a list of jobs from the specified namespace that match the statuses defined in the JOB_STATUSES environment variable.
// It returns a slice of job descriptions and an error if any occurs.
//
// Parameters:
// - clientset: A Kubernetes clientset to interact with the Kubernetes API.
// - namespace: The namespace from which to retrieve the jobs.
// - log: A logger to log messages.
//
// Returns:
// - A slice of ContainerInfo, each representing a job description with namespace, pod name, and status.
// - An error if any occurs during the retrieval of jobs.
func GetJobs(clientset *kubernetes.Clientset, namespace string, log *logrus.Logger) ([]ContainerInfo, error) {
	statuses := strings.Split(strings.TrimSpace(utils.GetEnv("JOB_STATUSES", "Complete", log)), ",")
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utils.LogWithFields(logrus.ErrorLevel, []string{}, "Error retrieving jobs", err)
		return nil, err
	}

	var jobsList []ContainerInfo
	for _, job := range jobs.Items {
		for _, jobStatus := range job.Status.Conditions {
			if utils.Contains(statuses, string(jobStatus.Type)) {
				jobsList = append(jobsList, ContainerInfo{
					Namespace: job.Namespace,
					PodName:   job.Name,
					Status:    string(jobStatus.Type),
				})
			}
		}
	}
	return jobsList, nil
}

// DeleteJobs deletes the specified jobs from the given namespace and logs the actions taken.
//
// Parameters:
// - clientset: A Kubernetes clientset to interact with the Kubernetes API.
// - jobs: A slice of ContainerInfo, each representing a job description with namespace, pod name, and status.
// - log: A logger to log messages.
func DeleteJobs(clientset *kubernetes.Clientset, jobs []ContainerInfo, log *logrus.Logger) {
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job ContainerInfo) {
			defer wg.Done()
			propagationPolicy := metav1.DeletePropagationBackground
			err := clientset.BatchV1().Jobs(job.Namespace).Delete(context.Background(), job.PodName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
			if err != nil {
				utils.LogWithFields(logrus.ErrorLevel, []string{fmt.Sprintf("job:%s", job.PodName)}, "Failed to delete job", err)
			} else {
				metrics.JobsPruned.WithLabelValues(job.Namespace, job.Status).Add(1) // Increment the counter
				utils.LogWithFields(logrus.InfoLevel, []string{fmt.Sprintf("job:%s", job.PodName)}, "Successfully deleted job")
			}
		}(job)
	}
	wg.Wait()
}
