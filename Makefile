build:
	go build

clean:
	rm -rf dataworks-cli
	rm -rf bin/dataworks-cli-*

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/dataworks-cli_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/dataworks-cli_darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/dataworks-cli_linux-amd64

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/dataworks-cli_windows-amd64

build-all: clean build-macos build-linux build-windows

compress-linux:
	upx ./bin/dataworks-cli_linux*
