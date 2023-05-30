package waldo

import (
	"os"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

type TriggerOptions struct {
	GitCommit   string
	Help        bool
	RuleName    string
	UploadToken string
	Verbose     bool
	Version     bool
}

type TriggerAction struct {
	ioStreams   *lib.IOStreams
	options     *TriggerOptions
	runtimeInfo *lib.RuntimeInfo
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
	assetVersion := ta.detectAssetVersion()
	verbose := ta.detectVerbose()

	ad := tool.NewAgentDownloader(assetVersion, data.CLIPrefix, verbose, ta.ioStreams, ta.runtimeInfo)

	path, err := ad.Download()

	if err != nil {
		return err
	}

	defer ad.Cleanup()

	return ta.executeAgent(path, os.Args[1:])
}

//-----------------------------------------------------------------------------

func (ta *TriggerAction) detectAssetVersion() string {
	if version := os.Getenv("WALDO_CLI_ASSET_VERSION"); len(version) > 0 {
		return version
	}

	return "latest"
}

func (ta *TriggerAction) detectVerbose() bool {
	if verbose := os.Getenv("WALDO_CLI_VERBOSE"); verbose == "1" {
		return true
	}

	return false
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
