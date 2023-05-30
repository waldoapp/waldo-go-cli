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
	ioStreams   *lib.IOStreams
	options     *RemoveOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewRemoveAction(options *RemoveOptions, ioStreams *lib.IOStreams) *RemoveAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &RemoveAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (ra *RemoveAction) Perform() error {
	if err := data.ValidateRecipeName(ra.options.RecipeName); err != nil {
		return err
	}

	cfg, _, err := data.SetupConfiguration(data.CreateKindNever)

	if err != nil {
		return err
	}

	if err := cfg.RemoveRecipe(ra.options.RecipeName); err != nil {
		return err
	}

	ra.ioStreams.Printf("\nRecipe %q successfully removed!\n", ra.options.RecipeName)

	return nil
}
