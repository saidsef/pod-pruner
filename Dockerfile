# Build
FROM golang:1.22 AS build
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
COPY ./ ./
RUN go build pruner.go

# Application
FROM scratch
COPY --from=build /app/pruner /
CMD ["/pruner"]
