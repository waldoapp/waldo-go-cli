package waldo

import (
	"errors"
	"fmt"
	"sort"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
)

type AddOptions struct {
	AppName     string
	Platform    string
	RecipeName  string
	UploadToken string
	Verbose     bool
}

type AddAction struct {
	ioStreams      *lib.IOStreams
	options        *AddOptions
	promptReader   *lib.PromptReader
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewAddAction(options *AddOptions, ioStreams *lib.IOStreams, overrides map[string]string) *AddAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &AddAction{
		ioStreams:      ioStreams,
		options:        options,
		promptReader:   ioStreams.PromptReader(),
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (aa *AddAction) Perform() error {
	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return err
	}

	if recipe, _ := cfg.FindRecipe(aa.options.RecipeName); recipe != nil {
		return fmt.Errorf("Recipe already added: %q", aa.options.RecipeName)
	}

	results, err := aa.findBuildPaths(cfg.BasePath())

	if err != nil {
		return err
	}

	recipe, err := aa.makeRecipe(results)

	if err != nil {
		return err
	}

	if err := cfg.AddRecipe(recipe); err != nil {
		return err
	}

	aa.ioStreams.Printf("\nAdded recipe %q to Waldo configuration\n", recipe.Name)

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

func (aa *AddAction) makeRecipe(items []*tool.FoundBuildPath) (*data.Recipe, error) {
	var (
		item *tool.FoundBuildPath
		err  error
	)

	if len(items) > 1 {
		item = aa.askBuildPath(items)
	} else {
		item = items[0]
	}

	recipe := &data.Recipe{
		Name:        aa.options.RecipeName,
		AppName:     aa.options.AppName,                    // for now…
		Flavor:      data.BuildFlavor(aa.options.Platform), // ditto…
		UploadToken: aa.options.UploadToken,                // ditto…
		BasePath:    item.AbsPath}

	ios := aa.ioStreams
	verbose := aa.options.Verbose

	switch item.BuildTool {
	case tool.BuildToolCustom:
		recipe.CustomBuilder, err = tool.MakeCustomBuilder(item.AbsPath, item.RelPath, verbose, ios)

	case tool.BuildToolExpo:
		recipe.ExpoBuilder, err = tool.MakeExpoBuilder(item.AbsPath, item.RelPath, verbose, ios)

	case tool.BuildToolFlutter:
		recipe.FlutterBuilder, err = tool.MakeFlutterBuilder(item.AbsPath, item.RelPath, verbose, ios)

	case tool.BuildToolGradle:
		recipe.Flavor = data.BuildFlavorAndroid
		recipe.GradleBuilder, err = tool.MakeGradleBuilder(item.AbsPath, item.RelPath, verbose, ios)

	case tool.BuildToolReactNative:
		recipe.ReactNativeBuilder, err = tool.MakeReactNativeBuilder(item.AbsPath, item.RelPath, verbose, ios)

	case tool.BuildToolXcode:
		recipe.Flavor = data.BuildFlavorIos
		recipe.XcodeBuilder, err = tool.MakeXcodeBuilder(item.AbsPath, item.RelPath, verbose, ios)

	default:
		return nil, fmt.Errorf("Unknown build tool: %s", item.BuildTool.String())
	}

	if err != nil {
		return nil, err
	}

	return recipe, nil
}
