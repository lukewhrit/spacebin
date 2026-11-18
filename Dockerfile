# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /src

# Copy module files and download depdencies first to leverage Docker cache
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy rest of the source now so code changes don't invalidate the cache for dependencies
COPY . .

# Compile the binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags "-s -w" -o /spacebin ./cmd/spacebin

# Make empty directory for SQLite database on non-root
RUN mkdir /data

# Final stage
FROM gcr.io/distroless/base-debian12:nonroot

# Copy the compiled binary from the build stage
COPY --from=builder /spacebin /spacebin

# Copy the data directory and give it to the nonroot user
COPY --from=builder --chown=65532:65532 /data /data

WORKDIR /data

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s CMD ["/spacebin", "healthcheck"]

# Entrypoint
ENTRYPOINT ["/spacebin"]
