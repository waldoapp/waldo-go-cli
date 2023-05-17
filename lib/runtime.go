package lib

import (
	"runtime"
	"strings"
)

type RuntimeInfo struct {
	Arch     string
	Platform string
}

//-----------------------------------------------------------------------------

func DetectRuntimeInfo() *RuntimeInfo {
	return &RuntimeInfo{
		Arch:     detectArch(),
		Platform: detectPlatform()}
}

//-----------------------------------------------------------------------------

func detectArch() string {
	arch := runtime.GOARCH

	switch arch {
	case "amd64":
		return "x86_64"

	default:
		return arch
	}
}

func detectPlatform() string {
	platform := runtime.GOOS

	switch platform {
	case "darwin":
		return "macOS"

	default:
		return strings.Title(platform)
	}
}
