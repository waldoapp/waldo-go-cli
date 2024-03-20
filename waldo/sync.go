package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type SyncOptions struct {
	Clean       bool
	GitBranch   string
	GitCommit   string
	RecipeName  string
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

//--------------------------------------ˆ---------------------------------------

func (sa *SyncAction) Perform() error {
	err := sa.performBuild()

	if err != nil {
		return err
	}

	return sa.performUpload()
}

//--------------------------------------ˆ---------------------------------------

func (sa *SyncAction) performBuild() error {
	options := &BuildOptions{
		Clean:      sa.options.Clean,
		RecipeName: sa.options.RecipeName,
		Verbose:    sa.options.Verbose}

	return NewBuildAction(
		options,
		sa.ioStreams).Perform()
}

func (sa *SyncAction) performUpload() error {
	options := &UploadOptions{
		GitBranch:   sa.options.GitBranch,
		GitCommit:   sa.options.GitCommit,
		Target:      sa.options.RecipeName,
		VariantName: sa.options.VariantName,
		Verbose:     sa.options.Verbose}

	return NewUploadAction(
		options,
		sa.ioStreams).Perform()
}
