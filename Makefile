APP_NAME := monitorchi
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: build clean amd

build:
	go build $(LDFLAGS) -o $(APP_NAME) .

clean:
	rm -f bin/*

amd:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(APP_NAME) .
