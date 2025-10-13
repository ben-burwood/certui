FROM golang:latest AS build

# Set working directory to /app/backend for Go build context
WORKDIR /app/backend

# Copy go.mod and go.sum for dependency resolution
COPY backend/go.mod backend/go.sum ./

RUN go mod download

# Copy the backend source code
COPY backend/. .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

FROM gcr.io/distroless/static-debian11 AS release

WORKDIR /

COPY --from=build /entrypoint /entrypoint

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/entrypoint"]
