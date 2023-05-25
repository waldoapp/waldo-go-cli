package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type VersionOptions struct {
}

type VersionAction struct {
	ioStreams   *lib.IOStreams
	options     *VersionOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewVersionAction(options *VersionOptions, ioStreams *lib.IOStreams) *VersionAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &VersionAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (va *VersionAction) Perform() error {
	va.ioStreams.Printf("\n%s %s (%s/%s)\n", data.CLIName, data.CLIVersion, va.runtimeInfo.Platform, va.runtimeInfo.Arch)

	return nil
}
