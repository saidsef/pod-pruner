# Installation

## Requirements

- A Kubernetes cluster, and permission to create a ClusterRole and a ClusterRoleBinding in it
- `kubectl` with kustomize support
- Go at the version `go.mod` sets, to build from source

## Kubernetes

The manifests live under [`deployment/`](https://github.com/saidsef/pod-pruner/tree/main/deployment). Clone the repository, create the namespace and apply them with kustomize:

```sh
git clone https://github.com/saidsef/pod-pruner.git
cd pod-pruner
kubectl create namespace pod-pruner
kubectl apply -k deployment/
```

`deployment/kustomization.yml` puts every object in the `pod-pruner` namespace but does not create it, so create the namespace before you apply.

The base creates these objects, all named `pod-pruner`:

| Kind | Purpose |
|------|---------|
| ServiceAccount | The identity the pruner runs as |
| ClusterRole | `get`, `list` and `delete` on pods and on `batch` Jobs |
| ClusterRoleBinding | Grants the ClusterRole to the ServiceAccount across the cluster |
| Deployment | One replica of the pruner |

The Deployment sets example values for `NAMESPACES`, `CONTAINER_STATUSES` and `RESOURCES`, and leaves `DRY_RUN` unset, so the pruner starts in dry run. Edit the `env` block in [`deployment/base/deployment.yaml`](https://github.com/saidsef/pod-pruner/blob/main/deployment/base/deployment.yaml) to set your own values from [Configuration](./configuration.md).

The pod runs as a non-root user with a read-only root filesystem and no Linux capabilities, under the `RuntimeDefault` seccomp profile. [Metrics](./metrics.md#scraping) covers its `prometheus.io/scrape` annotations.

## Container image

CI publishes the image to `ghcr.io/saidsef/pod-pruner` for `linux/amd64` and `linux/arm64`. CI tags every build `vYYYY.MM`, for the year and month it ran. It also tags a build from `main` as `latest`. A build from any other branch gets the branch name as its tag, with each `/` turned into `-`.

The image starts from `scratch` and holds nothing but the static binary.

## Building from source

```sh
go build -o pod-pruner ./pruner/pruner.go
docker build -t pod-pruner .
```

The binary only runs inside a pod. It reads the in-cluster service account config, has no kubeconfig support, and exits with `Kubernetes config error` anywhere else.

## Log output

```sh
kubectl logs -n pod-pruner deploy/pod-pruner
```

The pruner writes one JSON object per log line. Near the top of the log, a `Starting pruner` line records the settings the pruner is running with:

```json
{"dryRun":true,"interval":"2m0s","level":"info","msg":"Starting pruner","namespaces":["*"],"resources":["PODS","JOBS"],"time":"2026-10-05T09:00:00Z"}
```

`*` in `namespaces` means the pruner sweeps every namespace. Every sweep writes one line per resource type and namespace. In dry run that line reads `Dry run mode. The following containers would be deleted`, and its `containers` field lists each pod as `namespace/name (reason)`. A sweep with nothing to delete writes `No containers to prune`. Jobs follow the same pattern with `jobs` in place of `containers`.

## Turning deletion on

Set `DRY_RUN` to `false` once the dry-run lines list only the pods and Jobs you expect to lose:

```yaml
env:
  - name: DRY_RUN
    value: 'false'
```

Changing the `env` block rolls out a new pod, which deletes what it finds from its first sweep. Each delete writes a `Successfully deleted` or `Failed to delete` line naming the pod or Job.

## Next steps

Choose the namespaces and statuses to prune in [Configuration](./configuration.md).
