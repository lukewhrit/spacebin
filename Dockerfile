FROM golang:1.22.4-alpine3.19

RUN mkdir /opt/spacebin

COPY . /opt/spacebin
WORKDIR /opt/spacebin

# We need GCC and other packages for sqlite3 support
RUN apk add --no-cache build-base

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags "-s -w" -o bin/spacebin -tags sqlite ./cmd/spacebin/main.go

# Run the generated binary
CMD ["/opt/spacebin/bin/spacebin"]
