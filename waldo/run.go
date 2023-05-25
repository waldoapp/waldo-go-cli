package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type RunOptions struct {
	Interactive bool
	Preview     bool
	Verbose     bool
}

type RunAction struct {
	ioStreams   *lib.IOStreams
	options     *RunOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewRunAction(options *RunOptions, ioStreams *lib.IOStreams) *RunAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &RunAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (sa *RunAction) Perform() error {
	return nil
}
