PREFIX?=/usr/local

BUILD_DIR?=.build
INSTALL_DIR?=$(PREFIX)/bin

.PHONY: all build clean install uninstall

all: clean build install

build: build_linux build_macos build_windows

build_linux: build_linux_arm64 build_linux_x86_64

build_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/waldo-linux-arm64 main.go

build_linux_x86_64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/waldo-linux-x86_64 main.go

build_macos: build_macos_arm64 build_macos_x86_64

build_macos_arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/waldo-macos-arm64 main.go

build_macos_x86_64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/waldo-macos-x86_64 main.go

build_windows: build_windows_arm64 build_windows_x86_64

build_windows_arm64:
	GOOS=windows GOARCH=arm64 go build -o $(BUILD_DIR)/waldo-windows-arm64.exe main.go

build_windows_x86_64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/waldo-windows-x86_64.exe main.go

clean:
	@ go clean -i
	@ rm -rfv $(BUILD_DIR)

install: build_macos_x86_64
	@ install -d $(INSTALL_DIR)
	@ install -Cv $(BUILD_DIR)/waldo-macos-arm64 $(INSTALL_DIR)/waldo

uninstall:
	@ rm -fv $(INSTALL_DIR)/waldo
