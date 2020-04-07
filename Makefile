build:
	go build ./cmd/...

run:
	go run ./cmd/...

clean:
	rm update-tag

all: build

test_internal:
	go test -v ./internal/...

test: test_internal

