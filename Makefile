OUT := bin/spirit

.PHONY: clean

all: spirit

spirit: clean
	@go mod download
	go build --ldflags "-s -w" -o $(OUT) ./cmd/spirit/main.go

clean:
	rm -rf bin/

run: spirit
	./bin/spirit

format:
	go fmt ./...

test:
	go test ./... -v -race
