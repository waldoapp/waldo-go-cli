package waldo

import (
	"fmt"
	"os"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type UploadOptions struct {
	AppID         string
	BuildPath     string
	GitBranch     string
	GitCommit     string
	LegacyHelp    bool
	LegacyVersion bool
	UploadToken   string
	VariantName   string
	Verbose       bool
}

type UploadAction struct {
	ioStreams   *lib.IOStreams
	options     *UploadOptions
	runtimeInfo *lib.RuntimeInfo

	appID       string
	buildPath   string
	uploadToken string
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
		ua.ioStreams.Printf("\n%v\n", data.FullVersion())

		return nil
	}

	if err := ua.processOptions(); err != nil {
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

	return ua.executeAgent(path, ua.makeAgentArgs())
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) detectAppID() (string, error) {
	appID := ua.options.AppID

	var err error

	if strings.HasPrefix(ua.uploadToken, "u-") {
		err = data.ValidateAppID(appID)
	} else if len(appID) > 0 {
		err = fmt.Errorf("Option %q not allowed with CI token", "--app_id")
	}

	if err != nil {
		return "", err
	}

	return appID, nil
}

func (ua *UploadAction) detectBuildPath() (string, error) {
	buildPath := ua.options.BuildPath

	if len(buildPath) == 0 {
		return "", fmt.Errorf("No build path specified")
	}

	return buildPath, nil
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

func (ua *UploadAction) detectUploadToken() (string, error) {
	uploadToken := ua.options.UploadToken

	if len(uploadToken) == 0 {
		uploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	if len(uploadToken) == 0 {
		profile, _, err := data.SetupProfile(data.CreateKindNever)

		if err == nil {
			uploadToken = profile.APIToken
		}
	}

	if err := data.ValidateUploadToken(uploadToken); err != nil {
		return "", err
	}

	return uploadToken, nil
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

func (ua *UploadAction) makeAgentArgs() []string {
	args := []string{"upload"}

	if len(ua.appID) > 0 {
		args = append(args, "--app_id", ua.appID)
	}

	if len(ua.options.GitBranch) > 0 {
		args = append(args, "--git_branch", ua.options.GitBranch)
	}

	if len(ua.options.GitCommit) > 0 {
		args = append(args, "--git_commit", ua.options.GitCommit)
	}

	if len(ua.uploadToken) > 0 {
		args = append(args, "--upload_token", ua.uploadToken)
	}

	if len(ua.options.VariantName) > 0 {
		args = append(args, "--variant_name", ua.options.VariantName)
	}

	if ua.options.Verbose {
		args = append(args, "--verbose")
	}

	if len(ua.buildPath) > 0 {
		args = append(args, ua.buildPath)
	}

	return args
}

func (ua *UploadAction) processOptions() error {
	var err error

	ua.buildPath, err = ua.detectBuildPath()

	if err != nil {
		return err
	}

	ua.uploadToken, err = ua.detectUploadToken()

	if err != nil {
		return err
	}

	ua.appID, err = ua.detectAppID()

	if err != nil {
		return err
	}

	return nil
}
