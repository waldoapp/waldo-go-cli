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
	options        *RunOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewRunAction(options *RunOptions, overrides map[string]string) *RunAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &RunAction{
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (sa *RunAction) Perform() error {
	return nil
}
