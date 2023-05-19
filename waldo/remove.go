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
	var (
		cfg *data.Configuration
		err error
	)

	err = data.ValidateRecipeName(ra.options.RecipeName)

	if err == nil {
		cfg, _, err = data.SetupConfiguration(false)
	}

	if err == nil {
		err = cfg.RemoveRecipe(ra.options.RecipeName)
	}

	if err == nil {
		ra.ioStreams.Printf("\nRecipe %q successfully removed!\n", ra.options.RecipeName)
	}

	return err
}
