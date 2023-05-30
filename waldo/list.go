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
				la.ioStreams.Printf("%16.16s: %s\n", "build root", relPath)
			} else {
				la.ioStreams.Printf("%16.16s: (none)\n", "build root")
			}

			token := recipe.UploadToken

			if len(token) > 0 {
				la.ioStreams.Printf("%16.16s: %s\n", "upload token", token)
			} else {
				la.ioStreams.Printf("%16.16s: (none)\n", "upload token")
			}

			summary := recipe.Summarize()

			if len(summary) > 0 {
				la.ioStreams.Printf("%16.16s: %s\n", "build options", summary)
			} else {
				la.ioStreams.Printf("%16.16s: (none)\n", "build options")
			}
		}

		if la.options.UserInfo {
			if ud != nil {
				if am, _ := ud.FindMetadata(recipe); am != nil {
					if len(am.BuildPath) > 0 {
						if buildTime := la.formatTime(lib.GetModificationTimeUTC(am.BuildPath)); len(buildTime) > 0 {
							la.ioStreams.Printf("%16.16s: %s\n", "last build", buildTime)
						} else {
							la.ioStreams.Printf("%16.16s: (unknown)\n", "last build")
						}
					}

					if uploadTime := la.formatTime(am.UploadTime); len(uploadTime) > 0 {
						if uploadToken := am.UploadToken; len(uploadToken) > 0 {
							la.ioStreams.Printf("%16.16s: %s to %s\n", "last upload", uploadTime, uploadToken)
						} else {
							la.ioStreams.Printf("%16.16s: %s\n", "last upload", uploadTime)
						}
					} else {
						la.ioStreams.Printf("%16.16s: (unknown)\n", "last upload")
					}
				} else {
					la.ioStreams.Printf("%16.16s: (unknown)\n", "last build")
					la.ioStreams.Printf("%16.16s: (unknown)\n", "last upload")
				}
			} else {
				la.ioStreams.Printf("%16.16s: (unknown)\n", "last build")
				la.ioStreams.Printf("%16.16s: (unknown)\n", "last upload")
			}
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
