package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type SessionOptions struct {
	Language  string
	Model     string
	OSVersion string
	Verbose   bool
}

type SessionAction struct {
	options        *SessionOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewSessionAction(options *SessionOptions, overrides map[string]string) *SessionAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &SessionAction{
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (sa *SessionAction) Perform() error {
	return nil
}
