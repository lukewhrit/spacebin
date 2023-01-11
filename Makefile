OUT := bin/spirit

.PHONY: clean

all: spirit

spirit:
	@go mod download

	@if [ "$(NO_SQLITE)" = "1" ]; then\
		go build --ldflags "-s -w" -o $(OUT) ./cmd/spirit/main.go;\
	else\
		go build --ldflags "-s -w" -tags sqlite -o $(OUT) ./cmd/spirit/main.go;\
	fi

run: spirit
	./bin/spirit

format:
	go fmt ./...

test:
	go test ./... -v -race
