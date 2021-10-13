.PHONY: all build clean

all: clean build

build: build_darwin build_linux build_windows

build_darwin: build_darwin_amd64 build_darwin_arm64

build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o bin/waldo-darwin-amd64 waldo.go

build_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/waldo-darwin-arm64 waldo.go

build_linux: build_linux_amd64 build_linux_arm64

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/waldo-linux-amd64 waldo.go

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/waldo-linux-arm64 waldo.go

build_windows: build_windows_amd64 build_windows_arm64

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/waldo-windows-amd64.exe waldo.go

build_windows_arm64:
	GOOS=windows GOARCH=arm64 go build -o bin/waldo-windows-arm64.exe waldo.go

clean:
	@ go clean
	@ rm -rfv bin/*
