package waldo

import (
	"path/filepath"
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

	la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", "RECIPE NAME", "PLATFORM", "APP NAME", "BUILD TOOL")

	for _, recipe := range cfg.Recipes {
		if len(recipe.Name) == 0 {
			continue
		}

		appName := recipe.AppName

		if len(appName) == 0 {
			appName = "(unknown)"
		}

		buildTool := recipe.BuildTool().String()

		la.ioStreams.Printf("%-16.16s  %-8.8s  %-24.24s  %s\n", recipe.Name, recipe.Platform, appName, buildTool)

		if la.options.LongFormat {
			buildRoot := "(none)"
			uploadToken := "(none)"
			buildOptions := "(none)"

			absPath := filepath.Join(cfg.BasePath(), recipe.BasePath)

			if relPath := lib.MakeRelativeToCWD(absPath); len(relPath) > 0 {
				buildRoot = relPath
			}

			if token := recipe.UploadToken; len(token) > 0 {
				uploadToken = token
			}

			if summary := recipe.Summarize(); len(summary) > 0 {
				buildOptions = summary
			}

			la.ioStreams.Printf("%16.16s: %s\n", "build root", buildRoot)
			la.ioStreams.Printf("%16.16s: %s\n", "upload token", uploadToken)
			la.ioStreams.Printf("%16.16s: %s\n", "build options", buildOptions)
		}

		if la.options.UserInfo {
			buildPath := "(unknown)"
			lastBuild := "(unknown)"
			lastUpload := "(unknown)"

			if ud != nil {
				if am, _ := ud.FindMetadata(recipe); am != nil {
					if len(am.BuildPath) > 0 {
						buildPath = am.BuildPath

						if buildTime := la.formatTime(lib.GetModificationTimeUTC(am.BuildPath)); len(buildTime) > 0 {
							lastBuild = buildTime
						}
					}

					if uploadTime := la.formatTime(am.UploadTime); len(uploadTime) > 0 {
						if uploadToken := am.UploadToken; len(uploadToken) > 0 {
							lastUpload = uploadTime + " to " + uploadToken
						} else {
							lastUpload = uploadTime
						}
					}
				}
			}

			la.ioStreams.Printf("%16.16s: %s\n", "build path", buildPath)
			la.ioStreams.Printf("%16.16s: %s\n", "last build", lastBuild)
			la.ioStreams.Printf("%16.16s: %s\n", "last upload", lastUpload)
		}
	}

	return nil
}

//-----------------------------------------------------------------------------

func (la *ListAction) formatTime(t time.Time) string {
	if t == (time.Time{}) {
		return ""
	}

	return t.Format(time.RFC3339)
}
