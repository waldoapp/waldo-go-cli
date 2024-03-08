package data

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func BriefVersion() string {
	return fmt.Sprintf("%s %s", CLIName, CLIVersion)
}

func FullVersion() string {
	runtimeInfo := lib.DetectRuntimeInfo()

	return fmt.Sprintf("%s %s (%s/%s)", CLIName, CLIVersion, runtimeInfo.Platform, runtimeInfo.Arch)
}
