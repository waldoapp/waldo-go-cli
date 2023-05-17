package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type RemoveOptions struct {
	RecipeName string
	Verbose    bool
}

type RemoveAction struct {
	ioStreams      *lib.IOStreams
	options        *RemoveOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewRemoveAction(options *RemoveOptions, ioStreams *lib.IOStreams, overrides map[string]string) *RemoveAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &RemoveAction{
		ioStreams:      ioStreams,
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (ra *RemoveAction) Perform() error {
	name := ra.options.RecipeName

	cfg, _, err := data.SetupConfiguration(false)

	if err == nil {
		err = cfg.RemoveRecipe(name)
	}

	if err != nil {
		return err
	}

	ra.ioStreams.Printf("\nRemoved recipe %q from Waldo configuration\n", name)

	return nil
}
