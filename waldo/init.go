package waldo

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type InitOptions struct {
	Verbose bool
}

type InitAction struct {
	ioStreams   *lib.IOStreams
	options     *InitOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewInitAction(options *InitOptions, ioStreams *lib.IOStreams) *InitAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &InitAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (ia *InitAction) Perform() error {
	cfg, created, err := data.SetupConfiguration(true)

	if err != nil {
		return err
	}

	if created {
		ia.ioStreams.Printf("\nInitialized empty Waldo configuration at %q\n", cfg.Path())
	} else {
		ia.ioStreams.Printf("\nReinitialized existing Waldo configuration at %q\n", cfg.Path())
	}

	return nil
}

//-----------------------------------------------------------------------------

// func (ia *InitAction) addRecipes() error {
// 	firstTime := true

// 	for {
// 		if !ia.askAddRecipe(firstTime) {
// 			break
// 		}

// 		firstTime = false
// 	}

// 	return nil
// }

// func (ia *InitAction) askAddRecipe(firstTime bool) bool {
// 	prompt := ""

// 	if firstTime {
// 		prompt = "Add a recipe"
// 	} else {
// 		prompt = "Add another recipe"
// 	}

// 	return ia.promptReader.PromptReadYN(prompt, false, false)
// }
