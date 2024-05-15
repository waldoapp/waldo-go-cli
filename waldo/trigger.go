package waldo

import (
	"os"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type TriggerOptions struct {
	GitCommit     string
	LegacyHelp    bool
	LegacyVersion bool
	RuleName      string
	UploadToken   string
	Verbose       bool
}

type TriggerAction struct {
	ioStreams   *lib.IOStreams
	options     *TriggerOptions
	runtimeInfo *lib.RuntimeInfo

	uploadToken string
}

//-----------------------------------------------------------------------------

func NewTriggerAction(options *TriggerOptions, ioStreams *lib.IOStreams) *TriggerAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &TriggerAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (ta *TriggerAction) Perform() error {
	if ta.options.LegacyVersion {
		ta.ioStreams.Printf("\n%v\n", data.FullVersion())

		return nil
	}

	if err := ta.processOptions(); err != nil {
		return err
	}

	ad := api.NewAgentDownloader(
		ta.detectDownloadAssetVersion(),
		data.CLIPrefix,
		ta.detectDownloadVerbose(),
		ta.ioStreams,
		ta.runtimeInfo)

	path, err := ad.Download()

	if err != nil {
		return err
	}

	defer ad.Cleanup()

	return ta.executeAgent(path, ta.makeAgentArgs())
}

//-----------------------------------------------------------------------------

func (ta *TriggerAction) detectDownloadAssetVersion() string {
	if version := os.Getenv("WALDO_CLI_ASSET_VERSION"); len(version) > 0 {
		return version
	}

	return "latest"
}

func (ta *TriggerAction) detectDownloadVerbose() bool {
	if verbose := os.Getenv("WALDO_CLI_VERBOSE"); verbose == "1" {
		return true
	}

	return false
}

func (ta *TriggerAction) detectUploadToken() (string, error) {
	uploadToken := ta.options.UploadToken

	var err error

	if len(uploadToken) == 0 {
		uploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	err = data.ValidateCIToken(uploadToken)

	if err != nil {
		return "", err
	}

	return uploadToken, nil
}

func (ta *TriggerAction) enrichEnvironment() lib.Environment {
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

func (ta *TriggerAction) executeAgent(path string, args []string) error {
	task := lib.NewTask(path, args...)

	task.Env = ta.enrichEnvironment()
	task.IOStreams = ta.ioStreams

	return task.Execute()
}

func (ta *TriggerAction) makeAgentArgs() []string {
	args := []string{"trigger"}

	if len(ta.options.GitCommit) > 0 {
		args = append(args, "--git_commit", ta.options.GitCommit)
	}

	if len(ta.options.RuleName) > 0 {
		args = append(args, "--rule_name", ta.options.RuleName)
	}

	if len(ta.uploadToken) > 0 {
		args = append(args, "--upload_token", ta.uploadToken)
	}

	if ta.options.Verbose {
		args = append(args, "--verbose")
	}

	return args
}

func (ta *TriggerAction) processOptions() error {
	var err error

	ta.uploadToken, err = ta.detectUploadToken()

	if err != nil {
		return err
	}

	return nil
}
