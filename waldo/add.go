package waldo

import (
	"fmt"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/expo"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/flutter"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/reactnative"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type AddOptions struct {
	RecipeName string
	Verbose    bool
}

type AddAction struct {
	ioStreams    *lib.IOStreams
	options      *AddOptions
	promptReader *lib.PromptReader
	runtimeInfo  *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewAddAction(options *AddOptions, ioStreams *lib.IOStreams) *AddAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &AddAction{
		ioStreams:    ioStreams,
		options:      options,
		promptReader: ioStreams.PromptReader(),
		runtimeInfo:  runtimeInfo}
}

//-----------------------------------------------------------------------------

func (aa *AddAction) Perform() error {
	if err := data.ValidateRecipeName(aa.options.RecipeName); err != nil {
		return err
	}

	cfg, created, err := data.SetupConfiguration(data.CreateKindIfNeeded)

	if err != nil {
		return err
	}

	if created {
		aa.ioStreams.Printf("\nInitialized empty Waldo configuration at %q\n", cfg.Path())
	}

	if recipe, _ := cfg.FindRecipe(aa.options.RecipeName); recipe != nil {
		return fmt.Errorf("Recipe already added: %q", aa.options.RecipeName)
	}

	verbose := aa.options.Verbose
	ios := aa.ioStreams

	buildPaths, err := tool.DetectBuildPaths(cfg.BasePath(), verbose, ios)

	if err != nil {
		return err
	}

	buildPath, err := tool.DetermineBuildPath(buildPaths, verbose, ios)

	if err != nil {
		return err
	}

	recipe, err := aa.makeRecipe(cfg, buildPath)

	if err != nil {
		return err
	}

	if !aa.confirmRecipe(cfg.BasePath(), recipe) {
		return fmt.Errorf("Recipe canceled: %q", aa.options.RecipeName)
	}

	if err = cfg.AddRecipe(recipe); err != nil {
		return err
	}

	aa.ioStreams.Printf("\nRecipe %q successfully added!\n", recipe.Name)

	return nil
}

//-----------------------------------------------------------------------------

func (aa *AddAction) confirmRecipe(basePath string, recipe *data.Recipe) bool {
	aa.ioStreams.Printf("\nNew recipe %q:\n\n", recipe.Name)

	appName := recipe.AppName

	if len(appName) == 0 {
		appName = "(unknown)"
	}

	buildTool := recipe.BuildTool().String()
	buildRoot := "(none)"
	uploadToken := "(none)"
	buildOptions := "(none)"

	absPath := filepath.Join(basePath, recipe.BasePath)

	if relPath := lib.MakeRelativeToCWD(absPath); len(relPath) > 0 {
		buildRoot = relPath
	}

	if token := recipe.UploadToken; len(token) > 0 {
		uploadToken = token
	}

	if summary := recipe.Summarize(); len(summary) > 0 {
		buildOptions = summary
	}

	aa.ioStreams.Printf("  Platform:      %s\n", recipe.Platform)
	aa.ioStreams.Printf("  App name:      %s\n", appName)
	aa.ioStreams.Printf("  Build tool:    %s\n", buildTool)
	aa.ioStreams.Printf("  Build root:    %s\n", buildRoot)
	aa.ioStreams.Printf("  Upload token:  %s\n", uploadToken)
	aa.ioStreams.Printf("  Build options: %s\n", buildOptions)

	return aa.promptReader.ReadYN(fmt.Sprintf("Add recipe %q", recipe.Name))
}

func (aa *AddAction) makeRecipe(cfg *data.Configuration, buildPath *tool.BuildPath) (*data.Recipe, error) {
	verbose := aa.options.Verbose
	ios := aa.ioStreams

	recipe := &data.Recipe{
		Name:     aa.options.RecipeName,
		BasePath: lib.MakeRelative(buildPath.AbsPath, cfg.BasePath())}

	switch buildPath.BuildTool {
	case tool.BuildToolExpo:
		builder, appName, platform, err := expo.MakeExpoBuilder(buildPath.AbsPath, buildPath.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.ExpoBuilder = builder

	case tool.BuildToolFlutter:
		builder, appName, platform, err := flutter.MakeFlutterBuilder(buildPath.AbsPath, buildPath.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.FlutterBuilder = builder

	case tool.BuildToolGradle:
		builder, appName, platform, err := gradle.MakeGradleBuilder(buildPath.AbsPath, buildPath.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.GradleBuilder = builder

	case tool.BuildToolReactNative:
		builder, appName, platform, err := reactnative.MakeReactNativeBuilder(buildPath.AbsPath, buildPath.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.ReactNativeBuilder = builder

	case tool.BuildToolXcode:
		builder, appName, platform, err := xcode.MakeXcodeBuilder(buildPath.AbsPath, buildPath.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.XcodeBuilder = builder

	default:
		return nil, fmt.Errorf("Unknown build tool: %s", buildPath.BuildTool.String())
	}

	return recipe, nil
}
