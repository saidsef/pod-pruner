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

// ContainerInfo represents the information of a container within a Kubernetes cluster.
type ContainerInfo struct {
	Namespace string // Namespace is the Kubernetes namespace in which the container resides.
	PodName   string // PodName is the name of the pod that contains the container.
	Status    string // Status is the current status of the container (e.g., Running, Terminated).
}
