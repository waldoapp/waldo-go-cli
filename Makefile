.PHONY: all build clean

all: clean build

build: build_linux build_macos build_windows

build_linux: build_linux_arm64 build_linux_x86_64

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/waldo-linux-arm64 waldo.go

build_linux_x86_64:
	GOOS=linux GOARCH=amd64 go build -o bin/waldo-linux-x86_64 waldo.go

build_macos: build_macos_arm64 build_macos_x86_64

build_macos_arm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/waldo-macos-arm64 waldo.go

build_macos_x86_64:
	GOOS=darwin GOARCH=amd64 go build -o bin/waldo-macos-x86_64 waldo.go

build_windows: build_windows_arm64 build_windows_x86_64

build_windows_arm64:
	GOOS=windows GOARCH=arm64 go build -o bin/waldo-windows-arm64.exe waldo.go

build_windows_x86_64:
	GOOS=windows GOARCH=amd64 go build -o bin/waldo-windows-x86_64.exe waldo.go

clean:
	@ go clean -i
	@ rm -rfv bin
