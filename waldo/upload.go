package waldo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type UploadOptions struct {
	BuildPath   string
	GitBranch   string
	GitCommit   string
	UploadToken string
	VariantName string
	Verbose     bool
}

type UploadAction struct {
	apiBuildEndpoint string
	apiErrorEndpoint string
	ciInfo           *lib.CIInfo
	gitInfo          *lib.GitInfo
	ioStreams        *lib.IOStreams
	options          *UploadOptions
	runtimeInfo      *lib.RuntimeInfo
	wrapperName      string
	wrapperVersion   string
}

//-----------------------------------------------------------------------------

func NewUploadAction(options *UploadOptions, ioStreams *lib.IOStreams, overrides map[string]string) *UploadAction {
	ciInfo := lib.DetectCIInfo(true)
	gitInfo := lib.InferGitInfo(ciInfo.SkipCount)
	runtimeInfo := lib.DetectRuntimeInfo()

	return &UploadAction{
		apiBuildEndpoint: overrides["apiBuildEndpoint"],
		apiErrorEndpoint: overrides["apiErrorEndpoint"],
		ciInfo:           ciInfo,
		gitInfo:          gitInfo,
		ioStreams:        ioStreams,
		options:          options,
		runtimeInfo:      runtimeInfo,
		wrapperName:      overrides["wrapperName"],
		wrapperVersion:   overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) Perform() error {
	err := ua.validateUploadToken()

	if err != nil {
		return err
	}

	absBuildPath, buildSuffix, flavor, err := ua.validateBuildPath()

	if err != nil {
		return err
	}

	workingPath := ua.makeWorkingPath()

	err = os.RemoveAll(workingPath)

	if err == nil {
		err = os.MkdirAll(workingPath, 0755)
	}

	if err != nil {
		return err
	}

	defer os.RemoveAll(workingPath)

	absBuildPayloadPath := ua.makeBuildPayloadPath(workingPath, absBuildPath, buildSuffix)

	// 	ua.displaySummary()

	err = ua.createBuildPayload(absBuildPath, absBuildPayloadPath, buildSuffix)

	if err == nil {
		err = ua.uploadBuildWithRetry(absBuildPayloadPath, buildSuffix, flavor)
	}

	if err != nil {
		ua.uploadErrorWithRetry(err, flavor)
	} else {
		ua.ioStreams.Printf("\nBuild ‘%s’ successfully uploaded to Waldo!\n", filepath.Base(ua.options.BuildPath))
	}

	return err
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) authorization() string {
	return fmt.Sprintf("Upload-Token %s", ua.options.UploadToken)
}

func (ua *UploadAction) buildContentType(buildSuffix string) string {
	switch buildSuffix {
	case "apk":
		return lib.BinaryContentType

	case "app":
		return lib.ZipContentType

	default:
		return ""
	}
}

func (ua *UploadAction) checkBuildStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	if status == 401 {
		return errors.New("Upload token is invalid or missing!")
	}

	if status < 200 || status > 299 {
		return fmt.Errorf("Unable to upload build to Waldo, HTTP status: %d", status)
	}

	return nil
}

func (ua *UploadAction) checkErrorStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return fmt.Errorf("Unable to upload error to Waldo, HTTP status: %d", status)
	}

	return nil
}

func (ua *UploadAction) createBuildPayload(buildPath, buildPayloadPath, buildSuffix string) error {
	parentPath := filepath.Dir(buildPath)
	buildName := filepath.Base(buildPath)

	switch buildSuffix {
	case "apk":
		if !lib.IsRegularFile(buildPath) {
			break
		}

		return nil

	case "app":
		if !lib.IsDirectory(buildPath) {
			break
		}

		return lib.ZipFolder(buildPayloadPath, parentPath, buildName)

	default:
		break
	}

	return fmt.Errorf("Unable to read build at ‘%s’", buildPath)
}

func (ua *UploadAction) errorContentType() string {
	return lib.JsonContentType
}

func (ua *UploadAction) makeBuildPayloadPath(workingPath, buildPath, buildSuffix string) string {
	buildName := filepath.Base(buildPath)

	switch buildSuffix {
	case "app":
		return filepath.Join(workingPath, buildName+".zip")

	default:
		return buildPath
	}
}

func (ua *UploadAction) makeBuildURL(flavor string) string {
	buildURL := ua.apiBuildEndpoint

	if len(buildURL) == 0 {
		buildURL = data.DefaultAPIBuildEndpoint
	}

	query := make(url.Values)

	lib.AddIfNotEmpty(&query, "agentName", data.AgentName)
	lib.AddIfNotEmpty(&query, "agentVersion", data.AgentVersion)
	lib.AddIfNotEmpty(&query, "arch", ua.runtimeInfo.Arch)
	lib.AddIfNotEmpty(&query, "ci", ua.ciInfo.Provider.String())
	lib.AddIfNotEmpty(&query, "ciGitBranch", ua.ciInfo.GitBranch)
	lib.AddIfNotEmpty(&query, "ciGitCommit", ua.ciInfo.GitCommit)
	lib.AddIfNotEmpty(&query, "flavor", flavor)
	lib.AddIfNotEmpty(&query, "gitAccess", ua.gitInfo.Access.String())
	lib.AddIfNotEmpty(&query, "gitBranch", ua.gitInfo.Branch)
	lib.AddIfNotEmpty(&query, "gitCommit", ua.gitInfo.Commit)
	lib.AddIfNotEmpty(&query, "platform", ua.runtimeInfo.Platform)
	lib.AddIfNotEmpty(&query, "userGitBranch", ua.options.GitBranch)
	lib.AddIfNotEmpty(&query, "userGitCommit", ua.options.GitCommit)
	lib.AddIfNotEmpty(&query, "variantName", ua.options.VariantName)
	lib.AddIfNotEmpty(&query, "wrapperName", ua.wrapperName)
	lib.AddIfNotEmpty(&query, "wrapperVersion", ua.wrapperVersion)

	buildURL += "?" + query.Encode()

	return buildURL
}

func (ua *UploadAction) makeErrorPayload(err error) string {
	payload := ""

	lib.AppendIfNotEmpty(&payload, "agentName", data.AgentName)
	lib.AppendIfNotEmpty(&payload, "agentVersion", data.AgentVersion)
	lib.AppendIfNotEmpty(&payload, "arch", ua.runtimeInfo.Arch)
	lib.AppendIfNotEmpty(&payload, "ci", ua.ciInfo.Provider.String())
	lib.AppendIfNotEmpty(&payload, "ciGitBranch", ua.ciInfo.GitBranch)
	lib.AppendIfNotEmpty(&payload, "ciGitCommit", ua.ciInfo.GitCommit)
	lib.AppendIfNotEmpty(&payload, "message", err.Error())
	lib.AppendIfNotEmpty(&payload, "platform", ua.runtimeInfo.Platform)
	lib.AppendIfNotEmpty(&payload, "wrapperName", ua.wrapperName)
	lib.AppendIfNotEmpty(&payload, "wrapperVersion", ua.wrapperVersion)

	payload = "{" + payload + "}"

	return payload
}

func (ua *UploadAction) makeErrorURL() string {
	errorURL := ua.apiErrorEndpoint

	if len(errorURL) == 0 {
		errorURL = data.DefaultAPIErrorEndpoint
	}

	return errorURL
}

func (ua *UploadAction) makeWorkingPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("waldo-upload-%d", os.Getpid()))
}

func (ua *UploadAction) uploadBuild(buildPayloadPath, buildSuffix, flavor string, retryAllowed bool) (bool, error) {
	ua.ioStreams.Printf("\nUploading build to Waldo…\n")

	url := ua.makeBuildURL(flavor)

	file, err := os.Open(buildPayloadPath)

	if err != nil {
		return false, fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	defer file.Close()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, file)

	if err != nil {
		return false, fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", ua.authorization())

	if contentType := ua.buildContentType(buildSuffix); len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	req.Header.Add("User-Agent", ua.userAgent(flavor))

	if ua.options.Verbose {
		lib.DumpRequest(ua.ioStreams, req, false)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return retryAllowed, fmt.Errorf("Unable to upload build to Waldo, error: %v, url: %s", err, url)
	}

	if ua.options.Verbose {
		lib.DumpResponse(ua.ioStreams, rsp, true)
	}

	defer rsp.Body.Close()

	return retryAllowed && lib.ShouldRetry(rsp), ua.checkBuildStatus(rsp)
}

func (ua *UploadAction) uploadBuildWithRetry(buildPayloadPath, buildSuffix, flavor string) error {
	for attempts := 1; attempts <= data.MaxPostAttempts; attempts++ {
		retry, err := ua.uploadBuild(buildPayloadPath, buildSuffix, flavor, attempts < data.MaxPostAttempts)

		if !retry || err == nil {
			return err
		}

		ua.ioStreams.EmitError(data.AgentPrefix, err)

		ua.ioStreams.Printf("\nFailed upload attempts: %d -- retrying…\n\n", attempts)
	}

	return nil
}

func (ua *UploadAction) uploadError(err error, flavor string, retryAllowed bool) (bool, error) {
	url := ua.makeErrorURL()
	body := ua.makeErrorPayload(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return false, fmt.Errorf("Unable to upload error to Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", ua.authorization())
	req.Header.Add("Content-Type", ua.errorContentType())
	req.Header.Add("User-Agent", ua.userAgent(flavor))

	// if ua.options.Verbose {
	// 	lib.DumpRequest(ua.ioStreams, req, true)
	// }

	rsp, err := client.Do(req)

	if err != nil {
		return retryAllowed, fmt.Errorf("Unable to upload error to Waldo, error: %v, url: %s", err, url)
	}

	// if ua.options.Verbose {
	// 	lib.DumpResponse(ua.ioStreams, rsp, true)
	// }

	defer rsp.Body.Close()

	return retryAllowed && lib.ShouldRetry(rsp), ua.checkErrorStatus(rsp)
}

func (ua *UploadAction) uploadErrorWithRetry(err error, flavor string) error {
	for attempts := 1; attempts <= data.MaxPostAttempts; attempts++ {
		retry, tmpErr := ua.uploadError(err, flavor, attempts < data.MaxPostAttempts)

		if !retry || tmpErr == nil {
			return tmpErr
		}

		ua.ioStreams.EmitError(data.AgentPrefix, err)

		ua.ioStreams.Printf("\nFailed upload error attempts: %d -- retrying…\n\n", attempts)
	}

	return nil
}

func (ua *UploadAction) userAgent(flavor string) string {
	ci := ua.ciInfo.Provider.String()

	if ci == "Unknown" {
		ci = "Go CLI"
	}

	version := ua.wrapperVersion

	if len(version) == 0 {
		version = data.AgentVersion
	}

	return fmt.Sprintf("Waldo %s/%s v%s", ci, flavor, version)
}

func (ua *UploadAction) validateBuildPath() (string, string, string, error) {
	if len(ua.options.BuildPath) == 0 {
		return "", "", "", errors.New("Empty build path")
	}

	buildPath, err := filepath.Abs(ua.options.BuildPath)

	if err != nil {
		return "", "", "", err
	}

	buildSuffix := strings.TrimPrefix(filepath.Ext(buildPath), ".")

	switch buildSuffix {
	case "apk":
		return buildPath, buildSuffix, "Android", nil

	case "app":
		return buildPath, buildSuffix, "iOS", nil

	default:
		return "", "", "", fmt.Errorf("File extension of build path at ‘%s’ is not recognized", buildPath)
	}
}

func (ua *UploadAction) validateUploadToken() error {
	if len(ua.options.UploadToken) == 0 {
		return errors.New("Empty upload token")
	}

	return nil
}
