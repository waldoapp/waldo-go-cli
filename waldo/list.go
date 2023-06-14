package waldo

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type ListOptions struct {
	LongFormat bool
	UserInfo   bool
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
	cfg, _, err := data.SetupConfiguration(data.CreateKindNever)

	if err != nil {
		return err
	}

	var ud *data.UserData

	if la.options.UserInfo {
		ud = data.SetupUserData(cfg)
	}

	recipes := la.sortRecipes(cfg.Recipes)

	la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", "RECIPE NAME", "PLATFORM", "APP NAME", "BUILD TOOL")

	for _, recipe := range recipes {
		name := recipe.Name

		if len(name) == 0 {
			continue
		}

		platform := recipe.Platform
		appName := la.formatString(recipe.AppName, "(unknown)")
		buildTool := recipe.BuildTool().String()

		la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", name, platform, appName, buildTool)

		if la.options.LongFormat {
			absPath := filepath.Join(cfg.BasePath(), recipe.BasePath)

			buildRoot := la.formatString(lib.MakeRelativeToCWD(absPath), "(none)")
			uploadToken := la.formatString(recipe.UploadToken, "(none)")
			buildOptions := la.formatString(recipe.Summarize(), "(none)")

			la.ioStreams.Printf("%16.16s: %s\n", "build root", buildRoot)
			la.ioStreams.Printf("%16.16s: %s\n", "upload token", uploadToken)
			la.ioStreams.Printf("%16.16s: %s\n", "build options", buildOptions)
		}

		if la.options.UserInfo {
			var (
				buildPath  string
				lastBuild  string
				lastUpload string
			)

			if ud != nil {
				if am, _ := ud.FindMetadata(recipe); am != nil {
					buildPath = am.BuildPath

					if len(buildPath) > 0 {
						lastBuild = la.formatTime(lib.GetModificationTimeUTC(buildPath))
					}

					lastUpload = la.formatTime(am.UploadTime)

					if uploadToken := am.UploadToken; len(lastUpload) > 0 && len(uploadToken) > 0 {
						lastUpload += " to " + uploadToken
					}
				}
			}

			buildPath = la.formatString(buildPath, "(unknown)")
			lastBuild = la.formatString(lastBuild, "(unknown)")
			lastUpload = la.formatString(lastUpload, "(unknown)")

			la.ioStreams.Printf("%16.16s: %s\n", "build path", buildPath)
			la.ioStreams.Printf("%16.16s: %s\n", "last build", lastBuild)
			la.ioStreams.Printf("%16.16s: %s\n", "last upload", lastUpload)
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

func (la *ListAction) formatString(value, defaultValue string) string {
	if len(value) > 0 {
		return value
	}

	return defaultValue
}

func (la *ListAction) formatTime(value time.Time) string {
	if value == (time.Time{}) {
		return ""
	}

	return value.Format(time.RFC3339)
}

func (la *ListAction) sortRecipes(recipes []*data.Recipe) []*data.Recipe {
	sort.Slice(recipes, func(i, j int) bool {
		return strings.ToLower(recipes[i].Name) < strings.ToLower(recipes[j].Name)
	})

	return recipes
}
