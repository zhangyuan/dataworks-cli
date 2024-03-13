.PHONY: build
build:
	go build

.PHONY: clean
clean:
	rm -rf dataworks-helper
	rm -rf bin/dataworks-helper-*

.PHONY: build-macos
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/dataworks-helper_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/dataworks-helper_darwin-arm64

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/dataworks-helper_linux-amd64

.PHONE: build-all
build-all: clean build-macos build-linux

.PHONE: compress
compress-linux:
	upx ./bin/dataworks-helper_linux*
