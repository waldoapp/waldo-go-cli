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
	defaultWrapperVersion  = "1.0.0"
)

var (
	waldoBuildPath        string
	waldoBuildPayloadPath string
	waldoBuildSuffix      string
	waldoFlavor           string
	waldoGitAccess        string
	waldoGitBranch        string
	waldoGitCommit        string
	waldoPlatform         string
	waldoUploadToken      string
	waldoVariantName      string
	waldoVerbose          bool
	waldoWorkingPath      string
)

func checkBuildPath() {
	if len(waldoBuildPath) == 0 {
		failUsage(fmt.Errorf("Missing required argument: ‘build-path’"))
	}
}

func checkUploadToken() {
	if len(waldoUploadToken) == 0 {
		waldoUploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	if len(waldoUploadToken) == 0 {
		failUsage(fmt.Errorf("Missing required option: ‘--upload_token’"))
	}
}

func displaySummary(uploader *waldo.Uploader) {
	fmt.Printf("\n")
	fmt.Printf("Build path:          %s\n", summarize(uploader.BuildPath()))
	fmt.Printf("Git branch:          %s\n", summarize(uploader.GitBranch()))
	fmt.Printf("Git commit:          %s\n", summarize(uploader.GitCommit()))
	fmt.Printf("Upload token:        %s\n", summarizeSecure(uploader.UploadToken()))
	fmt.Printf("Variant name:        %s\n", summarize(uploader.VariantName()))

	if waldoVerbose {
		fmt.Printf("\n")
		fmt.Printf("Build payload path:  %s\n", summarize(uploader.BuildPayloadPath()))
		fmt.Printf("Inferred git branch: %s\n", summarize(uploader.InferredGitBranch()))
		fmt.Printf("Inferred git commit: %s\n", summarize(uploader.InferredGitCommit()))
	}

	fmt.Printf("\n")
}

func displayUsage() {
	fmt.Printf(`
OVERVIEW: Upload build to Waldo

USAGE: waldo [options] <build-path>

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

func displayVersion() {
	fmt.Printf("%s %s / %s\n", defaultWrapperNameFull, defaultWrapperVersion, waldo.Version())
}

func fail(err error) {
	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("waldo: %v\n", err))

	os.Exit(1)
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

	return overrides
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fail(fmt.Errorf("Unhandled panic: %v", err))
		}
	}()

	displayVersion()

	parseArgs()

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

func parseArgs() {
	args := os.Args[1:]

	if len(args) == 0 {
		displayUsage()

		os.Exit(0)
	}

	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		switch arg {
		case "--help":
			displayUsage()

			os.Exit(0)

		case "--git_branch":
			if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
				failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
			}

			waldoGitBranch = args[0]

			args = args[1:]

		case "--git_commit":
			if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
				failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
			}

			waldoGitCommit = args[0]

			args = args[1:]

		case "--upload_token":
			if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
				failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
			}

			waldoUploadToken = args[0]

			args = args[1:]

		case "--variant_name":
			if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
				failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
			}

			waldoVariantName = args[0]

			args = args[1:]

		case "--verbose":
			waldoVerbose = true

		case "--version":
			os.Exit(0) // version already displayed

		default:
			if strings.HasPrefix(arg, "-") {
				failUsage(fmt.Errorf("Unknown option: ‘%s’", arg))
			}

			if len(waldoBuildPath) == 0 {
				waldoBuildPath = arg
			} else {
				failUsage(fmt.Errorf("Unknown argument: ‘%s’", arg))
			}
		}
	}
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
