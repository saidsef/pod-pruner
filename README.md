# Pod Pruner: Kubernetes Container Pruner

[![Go Report Card](https://goreportcard.com/badge/github.com/saidsef/pod-pruner)](https://goreportcard.com/report/github.com/saidsef/pod-pruner)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/saidsef/pod-pruner)
[![GoDoc](https://godoc.org/github.com/saidsef/pod-pruner?status.svg)](https://pkg.go.dev/github.com/saidsef/pod-pruner?tab=doc)
![GitHub release(latest by date)](https://img.shields.io/github/v/release/saidsef/pod-pruner)
![Commits](https://img.shields.io/github/commits-since/saidsef/pod-pruner/latest.svg)
![GitHub](https://img.shields.io/github/license/saidsef/pod-pruner)

This is a Kubernetes application written in Go (Golang) that periodically prunes containers based on their statuses, in the namespaces you name or across every namespace in the cluster. The application can operate in a dry-run mode, allowing you to see which containers would be deleted without actually removing them.

## What is the use case

This application efficiently manages Kubernetes environments by periodically removing unnecessary containers based on their statuses, from the namespaces you name or from every namespace in the cluster, thereby freeing up resources. It includes a dry-run mode for users to preview which containers would be pruned without executing the deletion. This optimises resource usage and ensures a cleaner, more manageable cluster, while providing metrics via the `/metrics` endpoint.

## Alternatives

This application was inspired by [pod-reaper](https://github.com/saidsef/pod-reaper/tree/master). If you need an alternative, I suggest using [pod-reaper](https://github.com/saidsef/pod-reaper/tree/master).

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [How It Works](#how-it-works)

## Prerequisites

- Go (version 1.22 or later)
- Kubernetes cluster
- Access to the Kubernetes API from within the cluster

## Installation

1. Clone the repository:
```bash
git clone https://github.com/saidsef/pod-pruner.git
```

2. Build the application:
```bash
go build -o pod-pruner pruner/pruner.go
```

3. Ensure that the application is packaged into a Docker image and pushed to a container registry if you plan to deploy it in a Kubernetes environment.

## Configuration

The application requires certain environment variables to be set:

- `DRY_RUN`: Set to `"true"` to enable dry-run mode (default is `"true"`).
- `RESOURCES`: A comma-separated list of Kubernetes resources (default is `"PODS"`)
- `NAMESPACES`: A comma-separated list of namespaces to monitor for containers to prune. Unset or empty means every namespace in the cluster, which needs the `ClusterRole` and `ClusterRoleBinding` in `deployment/base/`. Namespaces are logged as `*` when the pruner runs cluster-wide.
- `CONTAINER_STATUSES`: A comma-separated list of container status reasons to filter by (e.g., `Error,ContainerStatusUnknown,Unknown,Completed`). Required - there is no default, and an unset or empty value is an error. See [Container statuses](#container-statuses) for the values that can match.
- `JOB_STATUSES`: A comma-separated list of jobs statuses to filter by (default is `Complete`).
- `INTERVAL`: How often to sweep, as a Go duration such as `90s` or `2m` (default is `120s`). A value that is not a positive duration falls back to the default.

### Container statuses

Each entry is matched literally against a container's `state.waiting.reason` or `state.terminated.reason`. Init containers are matched too, but only where they failed - an init container that terminated with exit code 0 reports the reason `Completed` as part of a healthy pod starting up, so matching it would prune every pod that has one. These are free-form strings in the Kubernetes API, set by the kubelet and the container runtime, so an entry that neither of them emits is ignored without error.

Waiting reasons, all set by the kubelet:

`ContainerCreating`, `PodInitializing`, `CrashLoopBackOff`, `ErrImagePull`, `ImagePullBackOff`, `ImageInspectError`, `ErrImageNeverPull`, `InvalidImageName`, `CreateContainerConfigError`, `PreCreateHookError`, `CreateContainerError`, `PreStartHookError`, `PostStartHookError`, `RunContainerError`.

Terminated reasons:

| Reason | Emitted by |
| --- | --- |
| `Completed`, `Error`, `OOMKilled` | containerd, CRI-O |
| `ContainerStatusUnknown` | kubelet, when a container cannot be located |
| `Unknown`, `StartError` | containerd |
| `seccomp killed` | CRI-O |

Pod-level reasons such as `Evicted`, `DeadlineExceeded` and `NodeAffinity` live on `pod.status.reason`, not on a container status, and so never match. The same applies to health statuses reported by other tooling, for example Argo CD's `Degraded`.

Example of setting environment variables in a Kubernetes deployment spec:

```bash
kubectl apply -k deployment/ -n pod-pruner
```

## Usage

Once the application is deployed, it will sweep the configured namespaces immediately and then every `INTERVAL`, which defaults to `120s`. With `NAMESPACES` unset it sweeps every namespace in the cluster. It will log the containers that are eligible for pruning based on their statuses. If dry-run mode is disabled, it will proceed to delete the identified containers.

## How It Works

1. **Environment Variables**: The application retrieves configuration values from environment variables.
2. **Kubernetes Client**: It creates a Kubernetes client using in-cluster configuration to interact with the Kubernetes API.
3. **Container Monitoring**: On startup and then every `INTERVAL`, it checks the configured namespaces - or every namespace, where `NAMESPACES` is unset - for containers that are in the defined states (e.g., `Waiting`, `Terminated`).
4. **Pruning Logic**: If containers are found, it either logs the containers that would be deleted (in dry-run mode) or deletes them from the cluster.

## Metrics

This includes metrics to monitor the pruning activities. The following metrics are available:

- **Pods Pruned**: Total number of pods pruned, labelled by namespace.
- **Containers Pruned**: Total number of containers pruned, labelled by namespace.
- **Jobs Pruned**: Total number of jobs pruned, labelled by namespace.

The metrics are exposed at the `/metrics` endpoint and can be accessed via a Prometheus server.

## Source

Our latest and greatest source of *Reverse Geocoding* can be found on [GitHub]. [Fork us](https://github.com/saidsef/pod-pruner/fork)!

## Contributing

We would :heart: you to contribute by making a [pull request](https://github.com/saidsef/pod-pruner/pulls).

Please read the official [Contribution Guide](./CONTRIBUTING.md) for more information on how you can contribute.
