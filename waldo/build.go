package waldo

import (
	"fmt"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

type BuildOptions struct {
	Clean      bool
	RecipeName string
	Verbose    bool
}

type BuildAction struct {
	ioStreams   *lib.IOStreams
	options     *BuildOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewBuildAction(options *BuildOptions, ioStreams *lib.IOStreams) *BuildAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &BuildAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (ba *BuildAction) Perform() error {
	if len(ba.options.RecipeName) > 0 {
		if err := data.ValidateRecipeName(ba.options.RecipeName); err != nil {
			return err
		}
	}

	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return err
	}

	recipe, err := cfg.FindRecipe(ba.options.RecipeName)

	if err != nil {
		return err
	}

	if err = ba.buildRecipe(cfg, recipe); err != nil {
		return err
	}

	ba.ioStreams.Printf("\nRecipe %q successfully built!\n", recipe.Name)

	return nil
}

//-----------------------------------------------------------------------------

func (ba *BuildAction) buildRecipe(cfg *data.Configuration, r *data.Recipe) error {
	ba.ioStreams.Printf("\nBuilding recipe %q…\n", r.Name)

	ud, err := data.SetupUserData(cfg)

	if err != nil {
		return err
	}

	absBasePath := filepath.Join(cfg.BasePath(), r.BasePath)

	var am *tool.ArtifactMetadata

	switch r.BuildTool() {
	case tool.BuildToolCustom:
		am, err = r.CustomBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolExpo:
		am, err = r.ExpoBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolFlutter:
		am, err = r.FlutterBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolGradle:
		am, err = r.GradleBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolReactNative:
		am, err = r.ReactNativeBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	case tool.BuildToolXcode:
		am, err = r.XcodeBuilder.Build(absBasePath, ba.options.Clean, ba.options.Verbose, ba.ioStreams)

	default:
		return fmt.Errorf("Don’t know how to build this app!")
	}

	if err != nil {
		return err
	}

	return ud.AddMetadata(r, am)
}
