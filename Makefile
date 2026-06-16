.PHONY: build run test lint fmt clean

BIN := aurview

build:
	go build -o bin/$(BIN) ./cmd/aurview

run:
	go run ./cmd/aurview $(ARGS)

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -w cmd internal

clean:
	rm -rf bin
