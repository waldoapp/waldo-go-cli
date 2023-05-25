package waldo

import (
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type ListOptions struct {
	LongFormat bool
}

type ListAction struct {
	ioStreams   *lib.IOStreams
	options     *ListOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewListAction(options *ListOptions, ioStreams *lib.IOStreams) *ListAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &ListAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (la *ListAction) Perform() error {
	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return err
	}

	la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", "RECIPE NAME", "PLATFORM", "APP NAME", "BUILD TOOL")

	for _, recipe := range cfg.Recipes {
		if len(recipe.Name) == 0 {
			continue
		}

		appName := recipe.AppName

		if len(appName) == 0 {
			appName = "(unknown)"
		}

		flavor := recipe.Flavor

		if len(flavor) == 0 {
			flavor = "Unknown"
		}

		tool := recipe.BuildTool().String()

		la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", recipe.Name, flavor, appName, tool)

		if la.options.LongFormat {
			absPath := filepath.Join(cfg.BasePath(), recipe.BasePath)
			relPath := lib.MakeRelativeToCWD(absPath)

			if len(relPath) > 0 {
				la.ioStreams.Printf("%-16.16s  build root: %s\n", "", relPath)
			}

			token := recipe.UploadToken

			if len(token) > 0 {
				la.ioStreams.Printf("%-16.16s  upload token: %s\n", "", token)
			}

			summary := recipe.Summarize()

			if len(summary) > 0 {
				la.ioStreams.Printf("%-16.16s  %s\n", "", summary)
			}
		}
	}

	return nil
}
