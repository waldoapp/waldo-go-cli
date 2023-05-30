package waldo

import (
	"errors"
	"fmt"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

type AddOptions struct {
	AppName     string
	Platform    string
	RecipeName  string
	UploadToken string
	Verbose     bool
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

	if len(aa.options.UploadToken) > 0 {
		if err := data.ValidateUploadToken(aa.options.UploadToken); err != nil {
			return err
		}
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

	results, err := aa.findBuildPaths(cfg.BasePath())

	if err != nil {
		return err
	}

	recipe, err := aa.makeRecipe(cfg, results)

	if err != nil {
		return err
	}

	if err = cfg.AddRecipe(recipe); err != nil {
		return err
	}

	aa.ioStreams.Printf("\nRecipe %q successfully added!\n", recipe.Name)

	return nil
}

//-----------------------------------------------------------------------------

func (aa *AddAction) askBuildPath(items []*tool.FoundBuildPath) *tool.FoundBuildPath {
	sort.Slice(items, func(i, j int) bool {
		return items[i].RelPath > items[j].RelPath
	})

	maxLen := 0

	for _, item := range items {
		rpLen := len(item.RelPath)

		if maxLen < rpLen {
			maxLen = rpLen
		}
	}

	choices := lib.Map(items, func(item *tool.FoundBuildPath) string {
		return fmt.Sprintf("%-*s  (%s)", maxLen, item.RelPath, item.BuildTool.String())
	})

	idx, _ := aa.promptReader.ReadChoose(
		"Available build paths",
		choices,
		"Choose a build path",
		0,
		false)

	return items[idx]
}

func (aa *AddAction) decideAppName(appName string) string {
	if len(aa.options.AppName) > 0 {
		return aa.options.AppName
	}

	return appName
}

func (aa *AddAction) decideBuildFlavor(platform string) data.BuildFlavor {
	flavor := data.ParseBuildFlavor(platform)

	if flavor != data.BuildFlavorUnknown {
		return flavor
	}

	return data.ParseBuildFlavor(platform)
}

func (aa *AddAction) findBuildPaths(rootPath string) ([]*tool.FoundBuildPath, error) {
	bd := tool.NewBuildDetector(aa.options.Verbose, aa.ioStreams)

	results, err := bd.Detect(rootPath)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("No build paths found")
	}

	return results, nil
}

func (aa *AddAction) makeRecipe(cfg *data.Configuration, items []*tool.FoundBuildPath) (*data.Recipe, error) {
	var item *tool.FoundBuildPath

	if len(items) > 1 {
		item = aa.askBuildPath(items)
	} else {
		item = items[0]
	}

	recipe := &data.Recipe{
		Name:        aa.options.RecipeName,
		UploadToken: aa.options.UploadToken, // for now…
		BasePath:    lib.MakeRelative(item.AbsPath, cfg.BasePath())}

	verbose := aa.options.Verbose
	ios := aa.ioStreams

	switch item.BuildTool {
	case tool.BuildToolCustom:
		builder, appName, platform, err := tool.MakeCustomBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.CustomBuilder = builder

	case tool.BuildToolExpo:
		builder, appName, platform, err := tool.MakeExpoBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.ExpoBuilder = builder

	case tool.BuildToolFlutter:
		builder, appName, platform, err := tool.MakeFlutterBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.FlutterBuilder = builder

	case tool.BuildToolGradle:
		builder, appName, platform, err := tool.MakeGradleBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.GradleBuilder = builder

	case tool.BuildToolReactNative:
		builder, appName, platform, err := tool.MakeReactNativeBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.ReactNativeBuilder = builder

	case tool.BuildToolXcode:
		builder, appName, platform, err := tool.MakeXcodeBuilder(item.AbsPath, item.RelPath, verbose, ios)

		if err != nil {
			return nil, err
		}

		recipe.AppName = aa.decideAppName(appName)     // for now…
		recipe.Flavor = aa.decideBuildFlavor(platform) // ditto…
		recipe.XcodeBuilder = builder

	default:
		return nil, fmt.Errorf("Unknown build tool: %s", item.BuildTool.String())
	}

	return recipe, nil
}
