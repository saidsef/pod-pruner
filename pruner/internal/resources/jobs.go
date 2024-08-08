package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/saidsef/pod-pruner/pruner/utils"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetJobs retrieves a list of jobs from the specified namespace that match the statuses defined in the JOB_STATUSES environment variable.
// It returns a slice of job descriptions and an error if any occurs.
//
// Parameters:
// - clientset: A Kubernetes clientset used to interact with the Kubernetes API.
// - namespace: The namespace from which to retrieve the jobs.
//
// Returns:
// - A slice of strings containing descriptions of the jobs that match the specified statuses.
// - An error if there was an issue retrieving the jobs.
func GetJobs(clientset *kubernetes.Clientset, namespace string, log *logrus.Logger) ([]string, error) {
	statuses := strings.Split(utils.GetEnv("JOB_STATUSES", "Complete", log), ",")
	jobs, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var jobsList []string
	for _, job := range jobs.Items {
		for _, jobStatus := range job.Status.Conditions {
			if utils.Contains(statuses, strings.TrimSpace(fmt.Sprint(jobStatus.Type))) {
				jobsList = append(jobsList, fmt.Sprintf("%s/%s: %s", job.Namespace, job.Name, jobStatus.Type))
			}
		}
	}
	return jobsList, nil
}

// DeleteJobs deletes the specified jobs from the given namespace and logs the actions taken.
//
// Parameters:
// - clientset: A Kubernetes clientset used to interact with the Kubernetes API.
// - namespace: The namespace from which to delete the jobs.
// - jobs: A slice of job names in the format "namespace/jobName" to be deleted.
// - log: A logger used to log the actions taken during the deletion process.
//
// This function does not return any values.
func DeleteJobs(clientset *kubernetes.Clientset, namespace string, jobs []string, log *logrus.Logger) {
	for _, job := range jobs {
		jobParts := strings.Split(job, "/")
		if len(jobParts) != 2 {
			log.Errorf("Invalid job format: %s", job)
			continue
		}
		jobName := strings.TrimSpace(strings.Split(jobParts[1], ":")[0]) // Trim any whitespace from the job name
		propagationPolicy := metav1.DeletePropagationBackground
		err := clientset.BatchV1().Jobs(namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{PropagationPolicy: &propagationPolicy})
		if err != nil {
			log.Errorf("Failed to delete job %s: %v", jobName, err)
		} else {
			log.Infof("Successfully deleted job %s", jobName)
		}
	}
}
