package waldo

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

type UploadOptions struct {
	GitBranch   string
	GitCommit   string
	Help        bool
	Target      string
	UploadToken string
	VariantName string
	Verbose     bool
	Version     bool
}

type UploadAction struct {
	ioStreams   *lib.IOStreams
	options     *UploadOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewUploadAction(options *UploadOptions, ioStreams *lib.IOStreams) *UploadAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &UploadAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) Perform() error {
	ciMode := ua.detectCIMode()
	args := os.Args[1:]

	var (
		ud          *data.UserData
		recipe      *data.Recipe
		uploadToken string
	)

	if !ciMode {
		ud, recipe = ua.detectPersistedData()

		if ud != nil && recipe != nil {
			var err error

			args, uploadToken, err = ua.mungeArgs(ud, recipe)

			if err != nil {
				return err
			}
		}
	}

	ad := tool.NewAgentDownloader(
		ua.detectAssetVersion(),
		data.CLIPrefix,
		ua.detectVerbose(),
		ua.ioStreams,
		ua.runtimeInfo)

	path, err := ad.Download()

	if err != nil {
		return err
	}

	defer ad.Cleanup()

	if err := ua.executeAgent(path, args); err != nil {
		return err
	}

	if !ciMode && ud != nil && recipe != nil {
		ua.updateMetadata(ud, recipe, uploadToken)
	}

	return nil
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) detectAssetVersion() string {
	if version := os.Getenv("WALDO_CLI_ASSET_VERSION"); len(version) > 0 {
		return version
	}

	return "latest"
}

func (ua *UploadAction) detectCIMode() bool {
	if ciMode := os.Getenv("CI"); ciMode == "true" || ciMode == "1" {
		return true
	}

	return false
}

func (ua *UploadAction) detectPersistedData() (*data.UserData, *data.Recipe) {
	cfg, _, err := data.SetupConfiguration(false)

	if err != nil {
		return nil, nil
	}

	ud, err := data.SetupUserData(cfg)

	if err != nil {
		return nil, nil
	}

	recipe, err := cfg.FindRecipe(ua.options.Target)

	if err != nil {
		return nil, nil
	}

	return ud, recipe
}

func (ua *UploadAction) detectVerbose() bool {
	if verbose := os.Getenv("WALDO_CLI_VERBOSE"); verbose == "1" {
		return true
	}

	return false
}

func (ua *UploadAction) enrichEnvironment() lib.Environment {
	env := lib.CurrentEnvironment()

	//
	// If _both_ wrapper override environment variables are already set, then
	// do _not_ replace either one with the CLI name/version (if only one or
	// neither is set, then go ahead and override both):
	//
	wrapperName := os.Getenv("WALDO_WRAPPER_NAME_OVERRIDE")
	wrapperVersion := os.Getenv("WALDO_WRAPPER_VERSION_OVERRIDE")

	if len(wrapperName) == 0 || len(wrapperVersion) == 0 {
		env["WALDO_WRAPPER_NAME_OVERRIDE"] = data.CLIName
		env["WALDO_WRAPPER_VERSION_OVERRIDE"] = data.CLIVersion
	}

	return env
}

func (ua *UploadAction) executeAgent(path string, args []string) error {
	task := lib.NewTask(path, args...)

	task.Env = ua.enrichEnvironment()
	task.IOStreams = ua.ioStreams

	return task.Execute()
}

func (ua *UploadAction) makeArgs(buildPath, uploadToken string) []string {
	args := []string{"upload", buildPath, "--upload_token", uploadToken}

	if len(ua.options.GitBranch) > 0 {
		args = append(args, "--git_branch", ua.options.GitBranch)
	}

	if len(ua.options.GitCommit) > 0 {
		args = append(args, "--git_commit", ua.options.GitCommit)
	}

	if ua.options.Help {
		args = append(args, "--help")
	}

	if len(ua.options.VariantName) > 0 {
		args = append(args, "--variant_name", ua.options.VariantName)
	}

	if ua.options.Verbose {
		args = append(args, "--verbose")
	}

	if ua.options.Version {
		args = append(args, "--version")
	}

	return args
}

func (ua *UploadAction) mungeArgs(ud *data.UserData, recipe *data.Recipe) ([]string, string, error) {
	uploadToken := ua.options.UploadToken

	if len(uploadToken) == 0 {
		uploadToken = recipe.UploadToken
	}

	if len(uploadToken) == 0 {
		uploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	err := data.ValidateUploadToken(uploadToken)

	if err != nil {
		return nil, "", err
	}

	target := ua.options.Target

	//
	// If not a recipe name, do not bother mungeing args:
	//
	if len(target) > 0 && strings.Contains(target, ".") {
		return os.Args[1:], uploadToken, nil
	}

	am, err := ud.FindMetadata(recipe)

	if err != nil {
		return nil, "", err
	}

	if len(am.BuildPath) == 0 {
		return nil, "", fmt.Errorf("No build found for recipe %q", recipe.Name)
	}

	return ua.makeArgs(am.BuildPath, uploadToken), uploadToken, nil
}

func (ua *UploadAction) updateMetadata(ud *data.UserData, recipe *data.Recipe, uploadToken string) {
	am, _ := ud.FindMetadata(recipe)

	if am != nil {
		am.UploadTime = time.Now().UTC()
		am.UploadToken = uploadToken

		ud.MarkDirty()

		ud.Save() // don’t care if save fails
	}
}
