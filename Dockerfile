FROM golang:1.14.6-alpine3.12

RUN mkdir /opt/spacebin-curiosity

ADD . /opt/spacebin-curiosity
WORKDIR /opt/spacebin-curiosity

# We need GCC and other packages for sqlite3 support
RUN apk add --no-cache build-base

# Download dependencies
RUN go mod download

# Build the binary
RUN go build --ldflags "-s -w" -tags sqlite ./

# Run the generated binary
CMD ["/opt/spacebin-curiosity/curiosity"]
