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
	"sync"
	"time"

	"github.com/saidsef/pod-pruner/pruner/internal/metrics"
	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// deleteConcurrency caps how many job deletions are in flight at once.
const deleteConcurrency = 10

// GetJobs retrieves a list of jobs from the specified namespace that match the statuses defined in the JOB_STATUSES environment variable.
// It returns a slice of job descriptions and an error if any occurs.
//
// Parameters:
// - clientset: A Kubernetes client interface to interact with the Kubernetes API.
// - namespace: The namespace from which to retrieve the jobs. An empty namespace means every namespace.
// - log: A logger to log messages.
//
// Returns:
// - A slice of ContainerInfo, each representing a job description with namespace, pod name, and status.
// - An error if any occurs during the retrieval of jobs.
func GetJobs(clientset kubernetes.Interface, namespace string, log *logrus.Logger) ([]ContainerInfo, error) {
	statuses := utils.SplitAndTrim(utils.GetEnv("JOB_STATUSES", "Complete", log))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var allJobs []ContainerInfo
	listOptions := metav1.ListOptions{Limit: listPageSize}

	for {
		jobList, err := clientset.BatchV1().Jobs(namespace).List(ctx, listOptions)
		if err != nil {
			utils.LogWithFields(logrus.ErrorLevel, []string{}, "Error retrieving jobs", err)
			return nil, err
		}

		for _, job := range jobList.Items {
			for _, jobStatus := range job.Status.Conditions {
				if utils.Contains(statuses, string(jobStatus.Type)) {
					allJobs = append(allJobs, ContainerInfo{
						Namespace: job.Namespace,
						PodName:   job.Name,
						Status:    string(jobStatus.Type),
					})
				}
			}
		}

		if jobList.Continue == "" {
			break
		}
		listOptions.Continue = jobList.Continue
	}

	return allJobs, nil
}

// DeleteJobs deletes the specified jobs from the given namespace and logs the actions taken.
//
// Parameters:
// - clientset: A Kubernetes client interface to interact with the Kubernetes API.
// - jobs: A slice of ContainerInfo, each representing a job description with namespace, pod name, and status.
// - log: A logger to log messages.
func DeleteJobs(clientset kubernetes.Interface, jobs []ContainerInfo, log *logrus.Logger) {
	forEachBounded(jobs, deleteConcurrency, func(job ContainerInfo) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		propagationPolicy := metav1.DeletePropagationBackground
		fields := []string{fmt.Sprintf("job:%s", job.PodName), fmt.Sprintf("namespace:%s", job.Namespace)}

		if err := clientset.BatchV1().Jobs(job.Namespace).Delete(ctx, job.PodName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy}); err != nil {
			utils.LogWithFields(logrus.ErrorLevel, fields, "Failed to delete job", err)
			return
		}
		metrics.JobsPruned.WithLabelValues(job.Namespace, job.Status).Add(1) // Increment the counter
		utils.LogWithFields(logrus.InfoLevel, fields, "Successfully deleted job")
	})
}

// forEachBounded runs fn over every item, at most limit at a time, and returns
// once they have all finished. Every call queues behind one client-side rate
// limiter, so an unbounded fan-out only means the later goroutines spend their
// timeout waiting for a token rather than talking to the API server.
//
// Parameters:
// - items: The items to visit.
// - limit: The greatest number of concurrent calls to fn.
// - fn: The function to run for each item.
func forEachBounded(items []ContainerInfo, limit int, fn func(ContainerInfo)) {
	var wg sync.WaitGroup
	tokens := make(chan struct{}, limit)

	for _, item := range items {
		wg.Add(1)
		tokens <- struct{}{}
		go func(item ContainerInfo) {
			defer wg.Done()
			defer func() { <-tokens }()
			fn(item)
		}(item)
	}
	wg.Wait()
}
