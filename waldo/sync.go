package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
)

type SyncOptions struct {
}

type SyncAction struct {
	options        *SyncOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewSyncAction(options *SyncOptions, overrides map[string]string) *SyncAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &SyncAction{
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (sa *SyncAction) Perform() error {
	return nil
}
