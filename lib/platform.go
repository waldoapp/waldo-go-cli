package lib

import "strings"

type Platform string

const (
	PlatformAndroid = "Android"
	PlatformIos     = "iOS"
	PlatformLinux   = "Linux"
	PlatformMacOS   = "macOS"
	PlatformUnknown = "Unknown"
	PlatformWindows = "Windows"
)

//-----------------------------------------------------------------------------

func ParsePlatform(value string) Platform {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "android":
		return PlatformAndroid

	case "darwin", "macos":
		return PlatformMacOS

	case "ios":
		return PlatformIos

	case "linux":
		return PlatformLinux

	case "windows":
		return PlatformWindows

	default:
		return PlatformUnknown
	}
}
