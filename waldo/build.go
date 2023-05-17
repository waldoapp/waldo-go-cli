package waldo

import (
	"fmt"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
)

type BuildOptions struct {
	RecipeName string
	Verbose    bool
}

type BuildAction struct {
	ioStreams      *lib.IOStreams
	options        *BuildOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewBuildAction(options *BuildOptions, ioStreams *lib.IOStreams, overrides map[string]string) *BuildAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &BuildAction{
		ioStreams:      ioStreams,
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (ba *BuildAction) Perform() error {
	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return err
	}

	recipe, err := cfg.FindRecipe(ba.options.RecipeName)

	if err != nil {
		return err
	}

	err = ba.buildRecipe(cfg, recipe)

	if err != nil {
		return err
	}

	ba.ioStreams.Printf("\nBuilt recipe %q from Waldo configuration\n", recipe.Name)

	return nil
}

//-----------------------------------------------------------------------------

func (ba *BuildAction) buildRecipe(cfg *data.Configuration, recipe *data.Recipe) error {
	absBasePath := filepath.Join(cfg.BasePath(), recipe.BasePath)

	switch recipe.BuildTool() {
	case tool.BuildToolCustom:
		return recipe.CustomBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolExpo:
		return recipe.ExpoBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolFlutter:
		return recipe.FlutterBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolGradle:
		return recipe.GradleBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolReactNative:
		return recipe.ReactNativeBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolXcode:
		return recipe.XcodeBuilder.Build(absBasePath, ba.options.Verbose, ba.ioStreams)

	default:
		return fmt.Errorf("Don’t know how to build this app!")
	}
}
