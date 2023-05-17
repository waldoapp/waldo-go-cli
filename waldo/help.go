package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type HelpOptions struct {
}

type HelpAction struct {
	options        *HelpOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewHelpAction(options *HelpOptions, overrides map[string]string) *HelpAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &HelpAction{
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (sa *HelpAction) Perform() error {
	return nil
}
