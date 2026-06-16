.PHONY: build run test lint fmt clean aur-srcinfo aur-verifysource aur-build release-check release-source-archive

BIN := aurview
AUR_DIR := packaging/aur

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

release-check:
	gofmt -w .
	go test ./...
	go vet ./...
	go build ./...

release-source-archive:
	test -n "$(VERSION)"
	./scripts/release/source-archive.sh "$(VERSION)" "$(AUR_DIR)/aurview_$(VERSION)_source.tar.gz"

aur-srcinfo:
	cd $(AUR_DIR) && makepkg --printsrcinfo > .SRCINFO

aur-verifysource:
	cd $(AUR_DIR) && makepkg --verifysource

aur-build:
	cd $(AUR_DIR) && makepkg -f

clean:
	rm -rf bin
