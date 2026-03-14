APP_NAME := monitorchi

.PHONY: build clean amd

build:
	go build -o $(APP_NAME) .


clean:
	rm -f bin/*

amd:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) .

