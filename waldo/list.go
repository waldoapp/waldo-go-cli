package waldo

import (
	"errors"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type ListOptions struct {
	LongFormat bool
}

type ListAction struct {
	ioStreams      *lib.IOStreams
	options        *ListOptions
	runtimeInfo    *lib.RuntimeInfo
	wrapperName    string
	wrapperVersion string
}

//-----------------------------------------------------------------------------

func NewListAction(options *ListOptions, ioStreams *lib.IOStreams, overrides map[string]string) *ListAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &ListAction{
		ioStreams:      ioStreams,
		options:        options,
		runtimeInfo:    runtimeInfo,
		wrapperName:    overrides["wrapperName"],
		wrapperVersion: overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (la *ListAction) Perform() error {
	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return err
	}

	recipes := cfg.Recipes

	if len(recipes) == 0 {
		return errors.New("No recipes defined")
	}

	lf := la.options.LongFormat

	for _, recipe := range recipes {
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

		la.ioStreams.Printf("  %-16.16s  %-7.7s  %-24.24s  %s\n", recipe.Name, flavor, appName, tool)

		if lf {
			path := recipe.BasePath

			if len(path) > 0 {
				la.ioStreams.Printf("  %-16.16s  %s\n", "", path)
			}

			summary := recipe.Summarize()

			if len(summary) > 0 {
				la.ioStreams.Printf("  %-16.16s  %s\n", "", summary)
			}

			token := recipe.UploadToken

			if len(token) > 0 {
				la.ioStreams.Printf("  %-16.16s  %s\n", "", token)
			}
		}
	}

	return nil
}
