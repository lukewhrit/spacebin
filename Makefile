OUT := bin/spacebin
MIGRATIONS_DIR := internal/database/migrations

.PHONY: clean migrate-up migrate-down

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

migrate-up:
	@if [ -z "$(MIGRATIONS_DRIVER)" ]; then echo "MIGRATIONS_DRIVER must be set (postgres|mysql|sqlite)"; exit 1; fi
	@command -v migrate >/dev/null 2>&1 || { echo "golang-migrate CLI (migrate) is required on PATH"; exit 127; }
	migrate -path $(MIGRATIONS_DIR)/$(MIGRATIONS_DRIVER) -database "$(SPIRIT_CONNECTION_URI)" up

migrate-down:
	@if [ -z "$(MIGRATIONS_DRIVER)" ]; then echo "MIGRATIONS_DRIVER must be set (postgres|mysql|sqlite)"; exit 1; fi
	@command -v migrate >/dev/null 2>&1 || { echo "golang-migrate CLI (migrate) is required on PATH"; exit 127; }
	migrate -path $(MIGRATIONS_DIR)/$(MIGRATIONS_DRIVER) -database "$(SPIRIT_CONNECTION_URI)" down
