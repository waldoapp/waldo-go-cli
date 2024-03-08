package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type VersionOptions struct {
}

type VersionAction struct {
	ioStreams *lib.IOStreams
	options   *VersionOptions
}

//-----------------------------------------------------------------------------

func NewVersionAction(options *VersionOptions, ioStreams *lib.IOStreams) *VersionAction {
	return &VersionAction{
		ioStreams: ioStreams,
		options:   options}
}

//-----------------------------------------------------------------------------

func (va *VersionAction) Perform() error {
	va.ioStreams.Printf("\n%s\n", data.FullVersion())

	return nil
}
