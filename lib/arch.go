package lib

import "strings"

type Arch string

const (
	ArchArm64   = "arm64"
	ArchUnknown = "unknown"
	ArchX86_64  = "x86_64"
)

//-----------------------------------------------------------------------------

func ParseArch(value string) Arch {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "amd64", "x86_64":
		return ArchX86_64

	case "arm64":
		return ArchArm64

	default:
		return ArchUnknown
	}
}
