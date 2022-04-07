package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-lib"
)

const (
	defaultWrapperName     = "Go CLI"
	defaultWrapperNameFull = "Waldo Go CLI"
	defaultWrapperVersion  = "1.1.2"
)

var (
	waldoBuildPath        string
	waldoBuildPayloadPath string
	waldoBuildSuffix      string
	waldoCommand          string
	waldoFlavor           string
	waldoGitAccess        string
	waldoGitBranch        string
	waldoGitCommit        string
	waldoPlatform         string
	waldoRuleName         string
	waldoUploadToken      string
	waldoVariantName      string
	waldoVerbose          bool
	waldoWorkingPath      string
)

func checkBuildPath() {
	if len(waldoBuildPath) == 0 {
		failMissingArg("build-path")
	}
}

func checkUploadToken() {
	if len(waldoUploadToken) == 0 {
		waldoUploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	if len(waldoUploadToken) == 0 {
		failMissingOpt("--upload_token")
	}
}

func displaySummary(context interface{}) {
	switch {
	case isTrigger():
		t := context.(*waldo.Triggerer)

		fmt.Printf("\n")
		fmt.Printf("Rule name:           %s\n", summarize(t.RuleName()))
		fmt.Printf("Upload token:        %s\n", summarizeSecure(t.UploadToken()))
		fmt.Printf("\n")

	case isUpload():
		u := context.(*waldo.Uploader)

		fmt.Printf("\n")
		fmt.Printf("Build path:          %s\n", summarize(u.BuildPath()))
		fmt.Printf("Git branch:          %s\n", summarize(u.GitBranch()))
		fmt.Printf("Git commit:          %s\n", summarize(u.GitCommit()))
		fmt.Printf("Upload token:        %s\n", summarizeSecure(u.UploadToken()))
		fmt.Printf("Variant name:        %s\n", summarize(u.VariantName()))

		if waldoVerbose {
			fmt.Printf("\n")
			fmt.Printf("Build payload path:  %s\n", summarize(u.BuildPayloadPath()))
			fmt.Printf("CI git branch:       %s\n", summarize(u.CIGitBranch()))
			fmt.Printf("CI git commit:       %s\n", summarize(u.CIGitCommit()))
			fmt.Printf("CI provider:         %s\n", summarize(u.CIProvider()))
			fmt.Printf("Git access:          %s\n", summarize(u.GitAccess()))
			fmt.Printf("Inferred git branch: %s\n", summarize(u.InferredGitBranch()))
			fmt.Printf("Inferred git commit: %s\n", summarize(u.InferredGitCommit()))
		}

		fmt.Printf("\n")
	}
}

func displayUsage() {
	switch {
	case isTrigger():
		fmt.Printf(`
OVERVIEW: Trigger run on Waldo

USAGE: waldo trigger [options]

OPTIONS:

  --help                  Display available options and exit
  --rule_name <value>     Rule name
  --upload_token <value>  Upload token (overrides WALDO_UPLOAD_TOKEN)
  --verbose               Display extra verbiage
  --version               Display version and exit
`)

	case isUpload():
		fallthrough

	default:
		fmt.Printf(`
OVERVIEW: Upload build to Waldo

USAGE: waldo upload [options] <build-path>

OPTIONS:

  --git_branch <value>    Branch name for originating git commit
  --git_commit <value>    Hash of originating git commit
  --help                  Display available options and exit
  --upload_token <value>  Upload token (overrides WALDO_UPLOAD_TOKEN)
  --variant_name <value>  Variant name
  --verbose               Display extra verbiage
  --version               Display version and exit
`)
	}
}

func displayVersion() {
	fmt.Printf("%s %s / %s\n", defaultWrapperNameFull, defaultWrapperVersion, waldo.Version())
}

func fail(err error) {
	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("waldo: %v\n", err))

	os.Exit(1)
}

func failMissingArg(arg string) {
	failUsage(fmt.Errorf("Missing required argument: ‘%s’", arg))
}

func failMissingOpt(opt string) {
	failUsage(fmt.Errorf("Missing required option: ‘%s’", opt))
}

func failMissingOptValue(opt string) {
	failUsage(fmt.Errorf("Missing required value for option: ‘%s’", opt))
}

func failUnknownArg(arg string) {
	failUsage(fmt.Errorf("Unknown argument: ‘%s’", arg))
}

func failUnknownOpt(opt string) {
	failUsage(fmt.Errorf("Unknown option: ‘%s’", opt))
}

func failUsage(err error) {
	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("waldo: %v\n", err))

	displayUsage()

	os.Exit(1)
}

func getOverrides() map[string]string {
	wrapperName := os.Getenv("WALDO_WRAPPER_NAME_OVERRIDE")

	if len(wrapperName) == 0 {
		wrapperName = defaultWrapperName
	}

	wrapperVersion := os.Getenv("WALDO_WRAPPER_VERSION_OVERRIDE")

	if len(wrapperVersion) == 0 {
		wrapperVersion = defaultWrapperVersion
	}

	overrides := map[string]string{
		"wrapperName":    wrapperName,
		"wrapperVersion": wrapperVersion}

	if apiBuildEndpoint := os.Getenv("WALDO_API_BUILD_ENDPOINT_OVERRIDE"); len(apiBuildEndpoint) > 0 {
		overrides["apiBuildEndpoint"] = apiBuildEndpoint
	}

	if apiErrorEndpoint := os.Getenv("WALDO_API_ERROR_ENDPOINT_OVERRIDE"); len(apiErrorEndpoint) > 0 {
		overrides["apiErrorEndpoint"] = apiErrorEndpoint
	}

	if apiTriggerEndpoint := os.Getenv("WALDO_API_TRIGGER_ENDPOINT_OVERRIDE"); len(apiTriggerEndpoint) > 0 {
		overrides["apiTriggerEndpoint"] = apiTriggerEndpoint
	}

	return overrides
}

func isTrigger() bool {
	return waldoCommand == "trigger"
}

func isUpload() bool {
	return waldoCommand == "upload"
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fail(fmt.Errorf("Unhandled panic: %v", err))
		}
	}()

	displayVersion()

	parseArgs()

	switch {
	case isTrigger():
		performTrigger()

	case isUpload():
		performUpload()
	}
}

func parseArgs() {
	args := os.Args[1:]

	if len(args) == 0 {
		displayUsage()

		os.Exit(0)
	}

	waldoCommand, args = parseCommand(args)

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch arg {
		case "--help":
			displayUsage()

			os.Exit(0)

		case "--git_branch":
			if isUpload() {
				waldoGitBranch, args = parseOption(arg, args)
			} else {
				failUnknownOpt(arg)
			}

		case "--git_commit":
			if isUpload() {
				waldoGitCommit, args = parseOption(arg, args)
			} else {
				failUnknownOpt(arg)
			}

		case "--rule_name":
			if isTrigger() {
				waldoRuleName, args = parseOption(arg, args)
			} else {
				failUnknownOpt(arg)
			}

		case "--upload_token":
			waldoUploadToken, args = parseOption(arg, args)

		case "--variant_name":
			if isUpload() {
				waldoVariantName, args = parseOption(arg, args)
			} else {
				failUnknownOpt(arg)
			}

		case "--verbose":
			waldoVerbose = true

		case "--version":
			os.Exit(0) // version already displayed

		default:
			if strings.HasPrefix(arg, "-") {
				failUnknownOpt(arg)
			}

			if isUpload() && len(waldoBuildPath) == 0 {
				waldoBuildPath = arg
			} else {
				failUnknownArg(arg)
			}
		}
	}
}

func parseCommand(args []string) (string, []string) {
	switch args[0] {
	case "trigger", "upload":
		return args[0], args[1:]

	default:
		return "upload", args
	}
}

func parseOption(arg string, args []string) (string, []string) {
	if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
		failMissingOptValue(arg)
	}

	return args[0], args[1:]
}

func performTrigger() {
	checkUploadToken()

	triggerer := waldo.NewTriggerer(
		waldoUploadToken,
		waldoRuleName,
		waldoVerbose,
		getOverrides())

	if err := triggerer.Validate(); err != nil {
		fail(err)
	}

	displaySummary(triggerer)

	fmt.Printf("Triggering run on Waldo\n")

	if err := triggerer.Perform(); err != nil {
		fail(err)
	}

	fmt.Printf("Run successfully triggered on Waldo!\n")
}

func performUpload() {
	checkBuildPath()
	checkUploadToken()

	uploader := waldo.NewUploader(
		waldoBuildPath,
		waldoUploadToken,
		waldoVariantName,
		waldoGitCommit,
		waldoGitBranch,
		waldoVerbose,
		getOverrides())

	if err := uploader.Validate(); err != nil {
		fail(err)
	}

	displaySummary(uploader)

	fmt.Printf("Uploading build to Waldo\n")

	if err := uploader.Upload(); err != nil {
		fail(err)
	}

	fmt.Printf("\nBuild ‘%s’ successfully uploaded to Waldo!\n", filepath.Base(waldoBuildPath))
}

func summarize(value string) string {
	if len(value) > 0 {
		return fmt.Sprintf("‘%s’", value)
	} else {
		return "(none)"
	}
}

func summarizeSecure(value string) string {
	if len(value) == 0 {
		return "(none)"
	}

	if !waldoVerbose {
		prefixLen := len(value)

		if prefixLen > 6 {
			prefixLen = 6
		}

		prefix := value[0:prefixLen]
		suffixLen := len(value) - len(prefix)
		secure := "********************************"

		value = prefix + secure[0:suffixLen]
	}

	return fmt.Sprintf("‘%s’", value)
}
