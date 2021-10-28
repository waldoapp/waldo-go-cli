package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

const (
	waldoAgentName        = "Waldo"
	waldoAgentVersion     = "0.3.0"
	waldoAPIBuildEndpoint = "https://api.waldo.io/versions"
	waldoAPIErrorEndpoint = "https://api.waldo.io/uploadError"
	waldoWrapperName      = "Go CLI"
	waldoWrapperVersion   = waldoAgentVersion
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
	waldoUserGitBranch    string
	waldoUserGitCommit    string
	waldoVariantName      string
	waldoVerbose          bool
	waldoWorkingPath      string
)

func addIfNotEmpty(query *url.Values, key string, value string) {
	if len(value) > 0 {
		query.Add(key, value)
	}
}

func checkBuildPath() {
	if len(waldoBuildPath) == 0 {
		failUsage(fmt.Errorf("Missing required argument: ‘build-path’"))
	}

	var err error

	waldoBuildPath, err = filepath.Abs(waldoBuildPath)

	if err != nil {
		fail(err)
	}

	waldoBuildSuffix = filepath.Ext(waldoBuildPath)

	if strings.HasPrefix(waldoBuildSuffix, ".") {
		waldoBuildSuffix = waldoBuildSuffix[1:]
	}

	switch waldoBuildSuffix {
	case "apk":
		waldoFlavor = "Android"

	case "app", "ipa":
		waldoFlavor = "iOS"

	default:
		fail(fmt.Errorf("File extension of build at ‘%s’ is not recognized", waldoBuildPath))
	}
}

func checkBuildStatus(resp *http.Response) {
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fail(err)
	}

	bodyString := string(body)

	statusRegex := regexp.MustCompile(`"status":([0-9]+)`)
	statusMatches := statusRegex.FindStringSubmatch(bodyString)

	if len(statusMatches) > 0 { // status is numeric _only_ on failure
		var status = 0

		status, err = strconv.Atoi(statusMatches[1])

		if err != nil {
			fail(err)
		}

		if status == 401 {
			fail(fmt.Errorf("Upload token is invalid or missing!"))
		}

		if status < 200 || status > 299 {
			fail(fmt.Errorf("Unable to upload build to Waldo, HTTP status: %d", status))
		}
	}
}

func checkGit() {
	if !isGitInstalled() {
		waldoGitAccess = "noGitCommandFound"
	} else if !hasGitRepository() {
		waldoGitAccess = "notGitRepository"
	} else {
		waldoGitAccess = "ok"
		waldoGitCommit = getGitCommit() // MUST do first
		waldoGitBranch = getGitBranch()
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

func checkVariantName() {
}

func createBuildPayload() {
	parentPath := filepath.Dir(waldoBuildPath)
	buildName := filepath.Base(waldoBuildPath)

	switch waldoBuildSuffix {
	case "app":
		if !isDir(waldoBuildPath) {
			fail(fmt.Errorf("Unable to read build at ‘%s’", waldoBuildPath))
		}

		waldoBuildPayloadPath = filepath.Join(waldoWorkingPath, buildName+".zip")

		err := zipDir(waldoBuildPayloadPath, parentPath, buildName)

		if err != nil {
			fail(err)
		}

	default:
		if !isRegular(waldoBuildPath) {
			fail(fmt.Errorf("Unable to read build at ‘%s’", waldoBuildPath))
		}

		waldoBuildPayloadPath = waldoBuildPath
	}
}

func createWorkingPath() {
	var err error

	waldoWorkingPath, err = os.MkdirTemp("", "WaldoGoCLI-*")

	if err != nil {
		fail(err)
	}

	os.RemoveAll(waldoWorkingPath)
	os.MkdirAll(waldoWorkingPath, 0755)
}

func deleteWorkingPath() {
	if len(waldoWorkingPath) > 0 {
		os.RemoveAll(waldoWorkingPath)
	}
}

func displaySummary() {
	fmt.Printf("\n")
	fmt.Printf("Build path:          %s\n", summarize(waldoBuildPath))
	fmt.Printf("Git branch:          %s\n", summarize(waldoUserGitBranch))
	fmt.Printf("Git commit:          %s\n", summarize(waldoUserGitCommit))
	fmt.Printf("Upload token:        %s\n", summarizeSecure(waldoUploadToken))
	fmt.Printf("Variant name:        %s\n", summarize(waldoVariantName))
	fmt.Printf("\n")

	if waldoVerbose {
		fmt.Printf("Build payload path:  %s\n", summarize(waldoBuildPayloadPath))
		fmt.Printf("Inferred git branch: %s\n", summarize(waldoGitBranch))
		fmt.Printf("Inferred git commit: %s\n", summarize(waldoGitCommit))
		fmt.Printf("\n")
	}
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
	waldoPlatform = getPlatform()

	fmt.Printf("%s %s %s (%s)\n", waldoAgentName, waldoWrapperName, waldoWrapperVersion, waldoPlatform)
}

func dumpRequest(req *http.Request, body bool) {
	if waldoVerbose {
		dump, err := httputil.DumpRequestOut(req, body)

		if err == nil {
			fmt.Printf("\n--- Request ---\n%s\n", dump)
		}
	}
}

func dumpResponse(resp *http.Response, body bool) {
	if waldoVerbose {
		dump, err := httputil.DumpResponse(resp, body)

		if err == nil {
			fmt.Printf("\n--- Response ---\n%s\n", dump)
		}
	}
}

func fail(err error) {
	message := fmt.Sprintf("waldo: %v", err)

	if len(waldoUploadToken) > 0 {
		if uploadError(err) {
			message += " -- Waldo team has been informed"
		}
	}

	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("%s\n", message))

	os.Exit(1)
}

func failUsage(err error) {
	if len(waldoUploadToken) > 0 {
		uploadError(err)
	}

	fmt.Printf("\n") // flush stdout

	os.Stderr.WriteString(fmt.Sprintf("waldo: %v\n", err))

	displayUsage()

	os.Exit(1)
}

func getArch() string {
	arch := runtime.GOARCH

	switch runtime.GOARCH {
	case "amd64":
		return "x86_64"

	default:
		return arch
	}
}

func getAuthorization() string {
	return fmt.Sprintf("Upload-Token %s", waldoUploadToken)
}

func getBuildContentType() string {
	switch waldoBuildSuffix {
	case "app":
		return "application/zip"

	default:
		return "application/octet-stream"
	}
}

func getCI() string {
	if len(os.Getenv("APPCENTER_BUILD_ID")) > 0 {
		return "App Center"
	}

	if len(os.Getenv("AGENT_ID")) > 0 {
		return "Azure DevOps"
	}

	if os.Getenv("BITRISE_IO") == "true" {
		return "Bitrise"
	}

	if len(os.Getenv("BUDDYBUILD_BUILD_ID")) > 0 {
		return "buddybuild"
	}

	if os.Getenv("CIRCLECI") == "true" {
		return "CircleCI"
	}

	if len(os.Getenv("CODEBUILD_BUILD_ID")) > 0 {
		return "CodeBuild"
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "GitHub Actions"
	}

	if len(os.Getenv("JENKINS_URL")) > 0 {
		return "Jenkins"
	}

	if len(os.Getenv("TEAMCITY_VERSION")) > 0 {
		return "TeamCity"
	}

	if os.Getenv("TRAVIS") == "true" {
		return "Travis CI"
	}

	if len(os.Getenv("CI_BUILD_ID")) > 0 {
		return "Xcode Cloud"
	}

	return ""
}

func getErrorContentType() string {
	return "application/json"
}

func getGitBranch() string {
	if len(waldoGitCommit) > 0 {
		name, _, err := run("git", "name-rev", "--refs=heads/*", "--name-only", waldoGitCommit)

		if err == nil && name != "HEAD" {
			return name
		}
	}

	name, _, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")

	if err == nil && name != "HEAD" {
		return name
	}

	return ""
}

func getGitCommit() string {
	skip := fmt.Sprintf("--skip=%d", getSkipCount())

	hash, _, err := run("git", "log", "--format=%H", skip, "-1")

	if err != nil {
		return ""
	}

	return hash
}

func getPlatform() string {
	platform := runtime.GOOS

	switch platform {
	case "darwin":
		return "macOS"

	default:
		return strings.Title(platform)
	}
}

func getSkipCount() int {
	if os.Getenv("GITHUB_ACTIONS") == "true" &&
		os.Getenv("GITHUB_EVENT_NAME") == "pull_request" {
		return 1
	}

	return 0
}

func getUserAgent() string {
	ci := getCI()

	if len(ci) == 0 {
		ci = waldoWrapperName
	}

	return fmt.Sprintf("%s %s/%s v%s", waldoAgentName, ci, waldoFlavor, waldoAgentVersion)
}

func hasGitRepository() bool {
	_, _, err := run("git", "rev-parse")

	return err == nil
}

func isDir(path string) bool {

	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return fi.Mode().IsDir()
}

func isGitInstalled() bool {
	var name string

	if runtime.GOOS == "windows" {
		name = "git.exe"
	} else {
		name = "git"
	}

	_, err := exec.LookPath(name)

	return err == nil
}

func isRegular(path string) bool {

	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return fi.Mode().IsRegular()
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
	checkGit()
	checkUploadToken()
	checkVariantName()

	createWorkingPath()

	defer deleteWorkingPath()

	createBuildPayload()

	displaySummary()

	uploadBuild()
}

func makeBuildURL() string {
	buildURL := os.Getenv("WALDO_API_BUILD_ENDPOINT_OVERRIDE")

	if len(buildURL) == 0 {
		buildURL = waldoAPIBuildEndpoint
	}

	query := make(url.Values)

	addIfNotEmpty(&query, "agentName", waldoAgentName)
	addIfNotEmpty(&query, "agentVersion", waldoAgentVersion)
	addIfNotEmpty(&query, "arch", getArch())
	addIfNotEmpty(&query, "ci", getCI())
	addIfNotEmpty(&query, "flavor", waldoFlavor)
	addIfNotEmpty(&query, "gitAccess", waldoGitAccess)
	addIfNotEmpty(&query, "gitBranch", waldoGitBranch)
	addIfNotEmpty(&query, "gitCommit", waldoGitCommit)
	addIfNotEmpty(&query, "platform", getPlatform())
	addIfNotEmpty(&query, "userGitBranch", waldoUserGitBranch)
	addIfNotEmpty(&query, "userGitCommit", waldoUserGitCommit)
	addIfNotEmpty(&query, "variantName", waldoVariantName)
	addIfNotEmpty(&query, "wrapperName", waldoWrapperName)
	addIfNotEmpty(&query, "wrapperVersion", waldoWrapperVersion)

	buildURL += "?" + query.Encode()

	return buildURL
}

func makeErrorPayload(err error) string {
	payload := fmt.Sprintf(`{"message":"%s"`, err.Error())

	if len(waldoAgentName) > 0 {
		payload += fmt.Sprintf(`,"agentName":"%s"`, waldoAgentName)
	}

	if len(waldoAgentVersion) > 0 {
		payload += fmt.Sprintf(`,"agentVersion":"%s"`, waldoAgentVersion)
	}

	arch := getArch()

	if len(arch) > 0 {
		payload += fmt.Sprintf(`,"arch":"%s"`, arch)
	}

	ci := getCI()

	if len(ci) > 0 {
		payload += fmt.Sprintf(`,"ci":"%s"`, ci)
	}

	platform := getPlatform()

	if len(platform) > 0 {
		payload += fmt.Sprintf(`,"platform":"%s"`, platform)
	}

	if len(waldoWrapperName) > 0 {
		payload += fmt.Sprintf(`,"wrapperName":"%s"`, waldoWrapperName)
	}

	if len(waldoWrapperVersion) > 0 {
		payload += fmt.Sprintf(`,"wrapperVersion":"%s"`, waldoWrapperVersion)
	}

	payload += "}"

	return payload

	return ""
}

func makeErrorURL() string {
	errorURL := os.Getenv("WALDO_API_ERROR_ENDPOINT_OVERRIDE")

	if len(errorURL) == 0 {
		errorURL = waldoAPIErrorEndpoint
	}

	return errorURL
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

			waldoUserGitBranch = args[0]

			args = args[1:]

		case "--git_commit":
			if len(args) == 0 || len(args[0]) == 0 || strings.HasPrefix(args[0], "-") {
				failUsage(fmt.Errorf("Missing required value for option: ‘%s’", arg))
			}

			waldoUserGitCommit = args[0]

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

func run(name string, args ...string) (string, string, error) {
	var (
		stderrBuffer bytes.Buffer
		stdoutBuffer bytes.Buffer
	)

	cmd := exec.Command(name, args...)

	cmd.Stderr = &stderrBuffer
	cmd.Stdout = &stdoutBuffer

	err := cmd.Run()

	stderr := strings.TrimRight(stderrBuffer.String(), "\n")
	stdout := strings.TrimRight(stdoutBuffer.String(), "\n")

	return stdout, stderr, err
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

func uploadBuild() {
	fmt.Printf("Uploading build to Waldo\n")

	buildName := filepath.Base(waldoBuildPath)

	url := makeBuildURL()

	file, err := os.Open(waldoBuildPayloadPath)

	if err != nil {
		fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
	}

	defer file.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, file)

	if err != nil {
		fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
	}

	req.Header.Add("Authorization", getAuthorization())
	req.Header.Add("Content-Type", getBuildContentType())
	req.Header.Add("User-Agent", getUserAgent())

	dumpRequest(req, false)

	resp, err := client.Do(req)

	if err != nil {
		fail(fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url))
	}

	dumpResponse(resp, true)

	checkBuildStatus(resp)

	defer resp.Body.Close()

	fmt.Printf("\nBuild ‘%s’ successfully uploaded to Waldo!\n", buildName)
}

func uploadError(err error) bool {
	url := makeErrorURL()
	body := makeErrorPayload(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return false
	}

	req.Header.Add("Authorization", getAuthorization())
	req.Header.Add("Content-Type", getErrorContentType())
	req.Header.Add("User-Agent", getUserAgent())

	// dumpRequest(req, true)

	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	defer resp.Body.Close()

	// dumpResponse(resp, true)

	return true
}

func zipDir(zipPath string, dirPath string, basePath string) error {
	err := os.Chdir(dirPath)

	if err != nil {
		return err
	}

	zipFile, err := os.Create(zipPath)

	if err != nil {
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	walker := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		zipEntry, err := zipWriter.Create(path)

		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)

		return err
	}

	err = filepath.WalkDir(basePath, walker)

	err2 := zipWriter.Close()

	if err != nil {
		return err
	}

	return err2
}
