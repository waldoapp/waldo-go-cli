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
	ioStreams   *lib.IOStreams
	options     *SessionOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewSessionAction(options *SessionOptions, ioStreams *lib.IOStreams) *SessionAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &SessionAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (sa *SessionAction) Perform() error {
	return nil
}
