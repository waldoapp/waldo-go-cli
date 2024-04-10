package waldo

import (
	"fmt"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/expo"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/flutter"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/ionic"
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
	userToken    string
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

	profile, _, err := data.SetupProfile(data.CreateKindNever)

	if err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	aa.userToken = profile.UserToken

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

	appID := recipe.AppID

	if len(appID) == 0 {
		appID = "(none)"
	}

	buildOptions := "(none)"
	buildRoot := "(none)"
	buildTool := recipe.BuildTool().String()

	absPath := filepath.Join(basePath, recipe.BasePath)

	if relPath := lib.MakeRelativeToCWD(absPath); len(relPath) > 0 {
		buildRoot = relPath
	}

	if summary := recipe.Summarize(); len(summary) > 0 {
		buildOptions = summary
	}

	aa.ioStreams.Printf("  App name:      %s\n", appName)
	aa.ioStreams.Printf("  App ID:        %s\n", appID)
	aa.ioStreams.Printf("  Platform:      %s\n", recipe.Platform)
	aa.ioStreams.Printf("  Build tool:    %s\n", buildTool)
	aa.ioStreams.Printf("  Build root:    %s\n", buildRoot)
	aa.ioStreams.Printf("  Build options: %s\n", buildOptions)

	return aa.promptReader.ReadYN(fmt.Sprintf("Add recipe %q", recipe.Name))
}

func (aa *AddAction) determineApp(platform lib.Platform) (string, string, error) {
	if aa.options.Verbose {
		aa.ioStreams.Printf("\nFetching %v apps for user token %q\n", platform, aa.userToken)
	}

	items, err := api.FetchApps(aa.userToken, platform, aa.options.Verbose, aa.ioStreams)

	if err != nil {
		return "", "", err
	}

	if platform != lib.PlatformUnknown {
		items = lib.CompactMap(items, func(item *tool.AppInfo) (*tool.AppInfo, bool) {
			return item, item.Platform == platform
		})
	}

	item, err := tool.DetermineApp(platform, items, aa.options.Verbose, aa.ioStreams)

	if err != nil {
		return "", "", err
	}

	return item.AppName, item.AppID, nil
}

func (aa *AddAction) makeRecipe(cfg *data.Configuration, buildPath *tool.BuildPath) (*data.Recipe, error) {
	verbose := aa.options.Verbose
	ios := aa.ioStreams

	recipe := &data.Recipe{
		Name:     aa.options.RecipeName,
		BasePath: lib.MakeRelative(buildPath.AbsPath, cfg.BasePath())}

	switch buildPath.BuildTool {
	case tool.BuildToolExpo:
		builder, appName, platform, err := expo.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.ExpoBuilder = builder

	case tool.BuildToolFlutter:
		builder, appName, platform, err := flutter.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.FlutterBuilder = builder

	case tool.BuildToolGradle:
		builder, appName, platform, err := gradle.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.GradleBuilder = builder

	case tool.BuildToolIonic:
		builder, appName, platform, err := ionic.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.IonicBuilder = builder

	case tool.BuildToolReactNative:
		builder, appName, platform, err := reactnative.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.ReactNativeBuilder = builder

	case tool.BuildToolXcode:
		builder, appName, platform, err := xcode.MakeBuilder(buildPath.AbsPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = appName
		recipe.Platform = platform
		recipe.XcodeBuilder = builder

	default:
		return nil, fmt.Errorf("Unknown build tool: %q", buildPath.BuildTool.String())
	}

	appName, appID, err := aa.determineApp(recipe.Platform)

	if err != nil {
		return nil, err
	}

	recipe.AppID = appID
	recipe.AppName = appName

	return recipe, nil
}
