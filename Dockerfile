# Build
FROM golang:1.24 AS build
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
COPY ./ ./
RUN go mod tidy && \
    go build -v -ldflags "-s -w" -trimpath -buildvcs -compiler gc -o ./pod-pruner ./pruner/pruner.go

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
