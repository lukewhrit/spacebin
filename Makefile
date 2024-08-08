OUT := bin/spacebin

.PHONY: clean

all: spacebin

spacebin: clean
	@go mod download
	go build --ldflags "-s -w" -o $(OUT) ./cmd/spacebin/main.go

clean:
	rm -rf bin/

run: spacebin
	./bin/spacebin

format:
	go fmt ./...

test:
	go test ./... -v -race

coverage:
	go test ./... -v -race -coverprofile=coverage.out
	go tool cover -html=coverage.out
