# Metrics

The pruner serves Prometheus metrics at `GET /metrics` on the port `PORT` sets, `8080` by default. The endpoint needs no authentication.

```sh
kubectl port-forward -n pod-pruner deploy/pod-pruner 8080
curl http://localhost:8080/metrics
```

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `pods_pruned_total` | counter | `namespace`, `state` | Pods the pruner deleted |
| `containers_pruned_total` | counter | `namespace`, `state` | The same series as `pods_pruned_total`, kept for dashboards that read it |
| `jobs_pruned_total` | counter | `namespace`, `state` | Jobs the pruner deleted |
| `go_*` | various | - | Go runtime statistics from the Prometheus client library |
| `process_*` | various | - | CPU, memory and file descriptor use of the pruner process |

`namespace` holds the namespace of the deleted object. `state` holds the container reason that matched for a pod, such as `Error`, or the condition type that matched for a Job, such as `Complete`.

A counter moves only when a delete succeeds. In dry run the pruner deletes nothing, so the endpoint carries no `*_pruned_total` series at all. A failed delete shows up in the log, not in a metric.

## Scraping

The Deployment in [`deployment/base/`](https://github.com/saidsef/pod-pruner/tree/main/deployment/base) sets `prometheus.io/scrape: "true"` and `prometheus.io/port: "8080"` on the pod. A Prometheus that discovers pods by annotation scrapes it with no further configuration.

The Prometheus Operator ignores those annotations. Under the operator, this PodMonitor selects the pod by label and scrapes its `metrics` port:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: pod-pruner
  namespace: pod-pruner
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: pod-pruner
  podMetricsEndpoints:
    - port: metrics
```

## Queries

```promql
# pods deleted per namespace over the last day
sum by (namespace) (increase(pods_pruned_total[1d]))

# pods deleted per reason over the last hour
sum by (state) (increase(pods_pruned_total[1h]))

# Jobs deleted per namespace over the last day
sum by (namespace) (increase(jobs_pruned_total[1d]))
```
