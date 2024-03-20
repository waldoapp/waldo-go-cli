package waldo

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type UploadOptions struct {
	AppID         string
	GitBranch     string
	GitCommit     string
	LegacyHelp    bool
	LegacyVersion bool
	Target        string
	UploadToken   string
	VariantName   string
	Verbose       bool
}

type UploadAction struct {
	ioStreams   *lib.IOStreams
	options     *UploadOptions
	runtimeInfo *lib.RuntimeInfo

	appID      string
	appToken   string
	buildPath  string
	recipeName string
	userToken  string
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
	if ua.options.LegacyVersion {
		ua.ioStreams.Printf("\n%s\n", data.FullVersion())

		return nil
	}

	ciMode := ua.detectCIMode()

	ud, recipe, err := ua.processOptions(ciMode)

	if err != nil {
		return err
	}

	ad := api.NewAgentDownloader(
		ua.detectDownloadAssetVersion(),
		data.CLIPrefix,
		ua.detectDownloadVerbose(),
		ua.ioStreams,
		ua.runtimeInfo)

	path, err := ad.Download()

	if err != nil {
		return err
	}

	defer ad.Cleanup()

	if err := ua.executeAgent(path, ua.makeAgentArgs()); err != nil {
		return err
	}

	if !ciMode && ud != nil && recipe != nil {
		ua.updateMetadata(ud, recipe)
	}

	return nil
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) detectAppID(required, allowed bool) (string, error) {
	appID := ua.options.AppID

	if len(appID) > 0 && !allowed {
		return "", fmt.Errorf("Option %q not allowed in this context", "--app_id")
	}

	if len(appID) > 0 || required {
		err := data.ValidateAppID(appID)

		if err != nil {
			return "", err
		}
	}

	return appID, nil
}

func (ua *UploadAction) detectAppToken(required, allowed bool) (string, error) {
	appToken := ua.options.UploadToken

	if len(appToken) > 0 && !allowed {
		return "", fmt.Errorf("Option %q not allowed in this context", "--upload_token")
	}

	if len(appToken) == 0 && allowed {
		appToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	if len(appToken) > 0 || required {
		err := data.ValidateAppToken(appToken)

		if err != nil {
			return "", err
		}
	}

	return appToken, nil
}

func (ua *UploadAction) detectCIMode() bool {
	if ciMode := os.Getenv("CI"); ciMode == "true" || ciMode == "1" {
		return true
	}

	return false
}

func (ua *UploadAction) detectDownloadAssetVersion() string {
	if version := os.Getenv("WALDO_CLI_ASSET_VERSION"); len(version) > 0 {
		return version
	}

	return "latest"
}

func (ua *UploadAction) detectDownloadVerbose() bool {
	if verbose := os.Getenv("WALDO_CLI_VERBOSE"); verbose == "1" {
		return true
	}

	return false
}

func (ua *UploadAction) detectPersistedData() (*data.UserData, *data.Recipe) {
	cfg, _, err := data.SetupConfiguration(data.CreateKindNever)

	if err != nil {
		return nil, nil
	}

	ud := data.SetupUserData(cfg)

	if ud == nil {
		return nil, nil
	}

	recipe, err := cfg.FindRecipe(ua.recipeName)

	if err != nil {
		return nil, nil
	}

	return ud, recipe
}

func (ua *UploadAction) detectRecipeName(required, allowed bool) (string, error) {
	recipeName := ua.options.Target

	if len(recipeName) > 0 && !allowed {
		return "", fmt.Errorf("Recipe name %q not allowed in this context", recipeName)
	}

	if len(recipeName) > 0 || required {
		err := data.ValidateRecipeName(recipeName)

		if err != nil {
			return "", err
		}
	}

	return recipeName, nil
}

func (ua *UploadAction) detectUserToken() (string, error) {
	profile, _, err := data.SetupProfile(data.CreateKindNever)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	return profile.UserToken, nil
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

func (ua *UploadAction) isRecipeTarget() bool {
	return !strings.Contains(ua.options.Target, ".")
}

func (ua *UploadAction) makeAgentArgs() []string {
	args := []string{"upload"}

	if len(ua.buildPath) > 0 {
		args = append(args, ua.buildPath)
	}

	if len(ua.appToken) == 0 {
		args = append(args, "--upload_token", ua.userToken)
		args = append(args, "--app_id", ua.appID)
	} else {
		args = append(args, "--upload_token", ua.appToken)

	}

	if len(ua.options.GitBranch) > 0 {
		args = append(args, "--git_branch", ua.options.GitBranch)
	}

	if len(ua.options.GitCommit) > 0 {
		args = append(args, "--git_commit", ua.options.GitCommit)
	}

	if len(ua.options.VariantName) > 0 {
		args = append(args, "--variant_name", ua.options.VariantName)
	}

	if ua.options.Verbose {
		args = append(args, "--verbose")
	}

	return args
}

func (ua *UploadAction) processOptions(ciMode bool) (*data.UserData, *data.Recipe, error) {
	isRecipe := ua.isRecipeTarget()

	var err error

	//
	// - An app token (formerly known as an “upload token”) is _allowed_
	//   only if the target is _not_ a recipe name.
	// - An app token is _required_ if and only if we are in CI Mode.
	//
	ua.appToken, err = ua.detectAppToken(ciMode, !isRecipe)

	if err != nil {
		return nil, nil, err
	}

	//
	// In CI mode, the target _must not_ be a recipe name:
	//
	if ciMode {
		_, err = ua.detectRecipeName(false, false)

		return nil, nil, err
	}

	//
	// Otherwise, the target _may_ be a recipe name iff an app token has _not_
	// been specified:
	//
	if isRecipe {
		ua.recipeName, err = ua.detectRecipeName(false, len(ua.appToken) == 0)

		if err != nil {
			return nil, nil, err
		}
	} else {
		ua.buildPath = ua.options.Target
	}

	//
	// If we now have a build path and a valid app token, we are done:
	//
	if len(ua.appToken) > 0 && len(ua.buildPath) > 0 {
		return nil, nil, err
	}

	//
	// Looks like we need a valid user token:
	//
	ua.userToken, err = ua.detectUserToken()

	if err != nil {
		return nil, nil, err
	}

	var (
		ud     *data.UserData
		recipe *data.Recipe
	)

	//
	// If no build path yet, try to find a recipe (possibly defaulted):
	//
	if len(ua.buildPath) == 0 {
		ud, recipe = ua.detectPersistedData()

		if ud != nil && recipe != nil {
			am, err := ud.FindMetadata(recipe)

			if err != nil {
				return nil, nil, err
			}

			if len(am.BuildPath) == 0 {
				return nil, nil, fmt.Errorf("No build found for recipe %q", recipe.Name)
			}

			ua.appID = recipe.AppID
			ua.buildPath = am.BuildPath
		}
	}

	//
	// If no build path yet, fail:
	//
	if len(ua.buildPath) == 0 {
		return nil, nil, fmt.Errorf("No build path specified")
	}

	//
	// If no app ID yet, require one:
	//
	if len(ua.appID) == 0 {
		ua.appID, err = ua.detectAppID(true, true)

		if err != nil {
			return nil, nil, err
		}
	}

	return ud, recipe, nil
}

func (ua *UploadAction) updateMetadata(ud *data.UserData, recipe *data.Recipe) {
	am, _ := ud.FindMetadata(recipe)

	if am != nil {
		if len(ua.appToken) == 0 {
			am.AppID = ua.appID
			am.UploadToken = ua.userToken
		} else {
			am.UploadToken = ua.appToken
		}

		am.UploadTime = time.Now().UTC()

		ud.MarkDirty()

		ud.Save() // don’t care if save fails
	}
}
