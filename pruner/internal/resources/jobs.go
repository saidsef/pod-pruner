package resources

import (
	"context"
	"fmt"
	"strings"
	"sync"

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
// - A slice of strings, each representing a job description in the format "namespace/jobName: jobStatus".
// - An error if any occurs during the retrieval of jobs.
func GetJobs(clientset *kubernetes.Clientset, namespace string, log *logrus.Logger) ([]string, error) {
	statuses := strings.Split(strings.TrimSpace(utils.GetEnv("JOB_STATUSES", "Complete", log)), ",")
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error retrieving jobs")
		return nil, err
	}

	var jobsList []string
	for _, job := range jobs.Items {
		for _, jobStatus := range job.Status.Conditions {
			if utils.Contains(statuses, string(jobStatus.Type)) {
				jobsList = append(jobsList, fmt.Sprintf("%s/%s: %s", job.Namespace, job.Name, jobStatus.Type))
			}
		}
	}
	return jobsList, nil
}

// DeleteJobs deletes the specified jobs from the given namespace and logs the actions taken.
//
// Parameters:
// - clientset: A Kubernetes clientset to interact with the Kubernetes API.
// - namespace: The namespace from which to delete the jobs.
// - jobs: A slice of strings, each representing a job description in the format "namespace/jobName: jobStatus".
// - log: A logger to log messages.
func DeleteJobs(clientset *kubernetes.Clientset, namespace string, jobs []string, log *logrus.Logger) {
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(job string) {
			defer wg.Done()
			jobParts := strings.Split(job, "/")
			if len(jobParts) != 2 {
				log.WithFields(logrus.Fields{
					"job": job,
				}).Error("Invalid job format")
				return
			}
			jobName := strings.TrimSpace(strings.Split(jobParts[1], ":")[0])
			propagationPolicy := metav1.DeletePropagationBackground
			err := clientset.BatchV1().Jobs(namespace).Delete(context.Background(), jobName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
			if err != nil {
				log.WithFields(logrus.Fields{
					"job":   jobName,
					"error": err,
				}).Error("Failed to delete job")
			} else {
				log.WithFields(logrus.Fields{
					"job": jobName,
				}).Info("Successfully deleted job")
			}
		}(job)
	}
	wg.Wait()
}
