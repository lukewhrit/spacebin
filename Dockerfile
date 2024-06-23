FROM golang:1.22.4-alpine3.19

RUN mkdir /opt/spirit

COPY . /opt/spirit
WORKDIR /opt/spirit

# We need GCC and other packages for sqlite3 support
RUN apk add --no-cache build-base

# Download dependencies
RUN go mod download

# Build the binary
RUN go build --ldflags "-s -w" -o bin/spirit -tags sqlite ./cmd/spirit/main.go

# Run the generated binary
CMD ["/opt/spirit/bin/spirit"]
