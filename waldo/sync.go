package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type SyncOptions struct {
	Clean       bool
	GitBranch   string
	GitCommit   string
	RecipeName  string
	UploadToken string
	VariantName string
	Verbose     bool
}

type SyncAction struct {
	ioStreams   *lib.IOStreams
	options     *SyncOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewSyncAction(options *SyncOptions, ioStreams *lib.IOStreams) *SyncAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &SyncAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (sa *SyncAction) Perform() error {
	return nil
}
