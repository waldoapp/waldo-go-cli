package waldo

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type VersionOptions struct {
}

type VersionAction struct {
	ioStreams      *lib.IOStreams
	options        *VersionOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewVersionAction(options *VersionOptions, ioStreams *lib.IOStreams, overrides map[string]string) *VersionAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &VersionAction{
		ioStreams:      ioStreams,
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (va *VersionAction) Perform() error {
	prefix := ""

	if len(va.wrapperName) > 0 && len(va.wrapperVersion) > 0 {
		prefix = fmt.Sprintf("%s %s / ", va.wrapperName, va.wrapperVersion)
	}

	va.ioStreams.Printf("\n%s%s %s (%s/%s)\n", prefix, data.AgentName, data.AgentVersion, va.runtimeInfo.Platform, va.runtimeInfo.Arch)

	return nil
}
