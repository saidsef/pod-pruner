# Build
FROM golang:1.26 AS build
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOFLAGS=-mod=readonly

# Dependencies resolve from the committed go.mod and go.sum, and cache as their
# own layer, so a source-only change does not refetch the module graph.
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -v -ldflags "-s -w" -trimpath -buildvcs -compiler gc -o ./pod-pruner ./pruner/pruner.go

# Application
FROM scratch

USER 1000

LABEL org.opencontainers.image.title="Pod Pruner"
LABEL org.opencontainers.image.description="Kubernetes Container Pruner"
LABEL org.opencontainers.image.source="https://github.com/saidsef/pod-pruner.git"
LABEL com.docker.extension.publisher-url="https://github.com/saidsef/pod-pruner.git"
LABEL com.docker.extension.categories="kubernetes,cleanup,pruner"

COPY --from=build /app/pod-pruner /
CMD ["/pod-pruner"]
