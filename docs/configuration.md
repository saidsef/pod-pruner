# Configuration

Every setting is an environment variable on the pruner container. The pruner trims the spaces around each entry in a comma-separated list and drops empty entries, so `PODS, JOBS` and `PODS,JOBS,` both work.

## Environment variables

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `DRY_RUN` | `true`, `false`, or another spelling Go's `strconv.ParseBool` accepts, such as `1`, `0` or `False` | `true` | With `false` the pruner deletes what each sweep finds. When `strconv.ParseBool` cannot read the value, the pruner logs an error and stays in dry run |
| `RESOURCES` | `PODS`, `JOBS` or both, comma-separated | `PODS` | The resource types to prune. Matching is case-sensitive, and the pruner ignores any other value |
| `NAMESPACES` | Namespace names, comma-separated | None, meaning every namespace | The namespaces to sweep. When unset or empty, the pruner sweeps every namespace and logs them as `*` |
| `CONTAINER_STATUSES` | Container reasons from [Container statuses](#container-statuses), comma-separated | None, required when `RESOURCES` includes `PODS` | The reasons that mark a pod for deletion. When unset or empty, every sweep logs `Error fetching containers` and skips each namespace, so the pruner prunes no Jobs either |
| `JOB_STATUSES` | Condition types from [Job statuses](#job-statuses), comma-separated | `Complete` | The condition types that mark a Job for deletion |
| `INTERVAL` | A positive Go duration, such as `90s`, `5m` or `1h30m` | `120s` | The time between sweeps. For a value that is not a positive duration, the pruner logs an error and uses the default |
| `PORT` | A TCP port number | `8080` | The port for the `/metrics` endpoint |

## Namespace scope

With `NAMESPACES` set, a sweep goes through the named namespaces one at a time. Without it, a sweep makes one cluster-wide list per resource type.

The manifests in [`deployment/base/`](https://github.com/saidsef/pod-pruner/tree/main/deployment/base) grant cluster-wide access through a ClusterRole and a ClusterRoleBinding, whichever way `NAMESPACES` is set. To keep the pruner out of the namespaces it does not sweep, replace the ClusterRoleBinding with a RoleBinding to the same ClusterRole in each namespace you list.

## Container statuses

The pruner compares each entry in `CONTAINER_STATUSES` with the `state.waiting.reason` and `state.terminated.reason` of every container in a pod. The match is exact and case-sensitive, and one match marks the whole pod for deletion.

Init containers count unless they exited with code 0. An init container that exits 0 reports `Completed` as a normal part of start-up, so matching it would delete every healthy pod that has an init container.

The kubelet and the container runtime set these reasons as free text, and the Kubernetes API has no fixed list of them. An entry that neither of them emits never matches, and the pruner does not warn about it.

The kubelet sets these waiting reasons:

| Reason | Meaning |
|--------|---------|
| `ContainerCreating` | The container does not exist yet, in a pod with no init containers |
| `PodInitializing` | The container does not exist yet, in a pod with init containers |
| `CrashLoopBackOff` | The container keeps exiting, and the kubelet is waiting out a growing back-off before the next restart |
| `ErrImagePull` | The image pull failed, for example on a missing tag, rejected credentials or an unreachable registry |
| `ImagePullBackOff` | An earlier pull failed, and the kubelet is waiting out a back-off before it tries again |
| `ImageInspectError` | The kubelet could not inspect the image |
| `ErrImageNeverPull` | `imagePullPolicy` is `Never` and the image is not on the node |
| `InvalidImageName` | The kubelet could not parse the image name |
| `CreateContainerConfigError` | The kubelet could not build the container config, for example when a referenced Secret or ConfigMap key is missing, or `runAsNonRoot` is set and the image runs as root |
| `PreCreateHookError` | The kubelet's internal step before container creation failed. That step applies CPU and memory manager pinning |
| `CreateContainerError` | The container runtime could not create the container |
| `PreStartHookError` | The kubelet's internal step before container start failed. That step registers the container with the CPU, memory and topology managers |
| `PostStartHookError` | The container's `postStart` lifecycle hook failed, and the kubelet killed the container |
| `RunContainerError` | The container runtime created the container but could not start it |
| `RestartingAllContainers` | The pod is restarting all its containers in place. The kubelet sets it only when the `RestartAllContainersOnContainerExits` feature gate is on |

`ContainerCreating` and `PodInitializing` are normal start-up states. Listing either one deletes healthy pods that happen to be starting when a sweep runs.

The container runtime or the kubelet sets these terminated reasons:

| Reason | Emitted by |
|--------|------------|
| `Completed`, `Error`, `OOMKilled` | containerd, CRI-O |
| `ContainerStatusUnknown` | kubelet, when it cannot locate a container |
| `Unknown`, `StartError` | containerd |
| `seccomp killed` | CRI-O |

Pod-level reasons such as `Evicted`, `DeadlineExceeded` and `NodeAffinity` live on `pod.status.reason`, not on a container status, so they never match. Health statuses from other tools, such as Argo CD's `Degraded`, never match either.

## Job statuses

The pruner compares each entry in `JOB_STATUSES` with the `type` of every condition on a Job. The Job controller sets these types, all defined in `k8s.io/api/batch/v1`:

| Type | Set when |
|------|----------|
| `Complete` | The Job finished successfully |
| `Failed` | The Job failed, for example by reaching its `backoffLimit` or `activeDeadlineSeconds` |
| `SuccessCriteriaMet` | The Job met its success criteria and is stopping its remaining pods, before `Complete` |
| `FailureTarget` | The Job is failing and is stopping its remaining pods, before `Failed` |
| `Suspended` | `spec.suspend` is true |

The pruner reads the condition type and ignores its `status`. A Job that was suspended and resumed keeps a `Suspended` condition with status `False`, and that condition still matches `Suspended`.
