package lib

import (
	"runtime"
)

type RuntimeInfo struct {
	Arch     Arch
	Platform Platform
}

//-----------------------------------------------------------------------------

func DetectRuntimeInfo() *RuntimeInfo {
	return &RuntimeInfo{
		Arch:     ParseArch(runtime.GOARCH),
		Platform: ParsePlatform(runtime.GOOS)}
}
