build:
	go build ./cmd/...

run:
	go run ./cmd/...

clean:
	rm update-tag

install:
	go install ./cmd/...

all: clean build

test_internal:
	go test -v ./internal/...

test: test_internal

