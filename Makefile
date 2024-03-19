build:
	go build

clean:
	rm -rf dataworks-helper
	rm -rf bin/dataworks-helper-*

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/dataworks-helper_darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/dataworks-helper_darwin-arm64

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/dataworks-helper_linux-amd64

build-windows:
    GOOS=windows GOARCH=amd64 go build -o bin/dataworks-helper_windows-amd64

build-all: clean build-macos build-linux build-windows

compress-linux:
	upx ./bin/dataworks-helper_linux*
