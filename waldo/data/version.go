package data

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func BriefVersion() string {
	return fmt.Sprintf("%v %v", CLIName, CLIVersion)
}

func FullVersion() string {
	runtimeInfo := lib.DetectRuntimeInfo()

	return fmt.Sprintf("%v %v (%v/%v)", CLIName, CLIVersion, runtimeInfo.Platform, runtimeInfo.Arch)
}
