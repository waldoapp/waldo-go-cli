package waldo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
	"github.com/waldoapp/waldo-go-cli/waldo/data/tool"
)

type UploadOptions struct {
	GitBranch   string
	GitCommit   string
	Target      string
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
	var (
		am               *tool.ArtifactMetadata
		flavor           data.BuildFlavor
		ud               *data.UserData
		recipe           *data.Recipe
		buildPayloadPath string
		uploadToken      string
		workingPath      string
		buildPath        string
		err              error
	)

	ud, recipe = ua.detectPersistedData()

	buildPath, flavor, err = ua.determineTarget(ud, recipe)

	if err == nil {
		uploadToken, err = ua.determineUploadToken(recipe)
	}

	if err == nil {
		workingPath, err = ua.prepareWorkingPath()
	}

	if err == nil {
		defer os.RemoveAll(workingPath)

		buildPayloadPath = ua.makeBuildPayloadPath(workingPath, buildPath, flavor)

		ua.displaySummary(uploadToken, buildPath, buildPayloadPath)

		err = ua.createBuildPayload(buildPath, buildPayloadPath, flavor)
	}

	if err == nil {
		err = ua.uploadBuildWithRetry(uploadToken, buildPayloadPath, flavor)
	}

	if err == nil && ud != nil && recipe != nil {
		am, _ = ud.FindMetadata(recipe)

		if am != nil {
			am.UploadTime = time.Now().UTC()
			am.UploadToken = uploadToken

			_ = ud.Save() // don’t care if save fails
		}
	}

	if err == nil {
		ua.ioStreams.Printf("\nBuild %q successfully uploaded to Waldo!\n", filepath.Base(ua.options.Target))
	} else {
		ua.uploadErrorWithRetry(err, uploadToken, flavor)
	}

	return err
}

//-----------------------------------------------------------------------------

func (ua *UploadAction) authorization(uploadToken string) string {
	return fmt.Sprintf("Upload-Token %s", uploadToken)
}

func (ua *UploadAction) buildContentType(flavor data.BuildFlavor) string {
	switch flavor {
	case data.BuildFlavorAndroid:
		return lib.BinaryContentType

	case data.BuildFlavorIos:
		return lib.ZipContentType

	default:
		return ""
	}
}

func (ua *UploadAction) checkBuildStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	switch {
	case status == 401:
		return errors.New("Upload token is invalid or missing!")

	case status < 200 || status > 299:
		return fmt.Errorf("Unable to upload build to Waldo, HTTP status: %d", status)

	default:
		return nil
	}
}

func (ua *UploadAction) checkErrorStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return fmt.Errorf("Unable to upload error to Waldo, HTTP status: %d", status)
	}

	return nil
}

func (ua *UploadAction) createBuildPayload(buildPath, buildPayloadPath string, flavor data.BuildFlavor) error {
	parentPath := filepath.Dir(buildPath)
	buildName := filepath.Base(buildPath)

	switch flavor {
	case data.BuildFlavorAndroid:
		if !lib.IsRegularFile(buildPath) {
			break
		}

		return nil

	case data.BuildFlavorIos:
		if !lib.IsDirectory(buildPath) {
			break
		}

		return lib.ZipFolder(buildPayloadPath, parentPath, buildName)

	default:
		break
	}

	return fmt.Errorf("Unable to read build at %q", buildPath)
}

func (ua *UploadAction) detectPersistedData() (*data.UserData, *data.Recipe) {
	var (
		cfg    *data.Configuration
		ud     *data.UserData
		recipe *data.Recipe
		err    error
	)

	cfg, _, err = data.SetupConfiguration(false)

	if err == nil {
		ud, err = data.SetupUserData(cfg)
	}

	if err == nil {
		recipe, err = cfg.FindRecipe(ua.options.Target)
	}

	if err == nil {
		return ud, recipe
	}

	return nil, nil
}

func (ua *UploadAction) determineTarget(ud *data.UserData, recipe *data.Recipe) (string, data.BuildFlavor, error) {
	var (
		am        *tool.ArtifactMetadata
		flavor    data.BuildFlavor
		buildPath string
		err       error
	)

	target := ua.options.Target

	if len(target) > 0 && strings.Contains(target, ".") {
		buildPath, err = filepath.Abs(target)
	} else if ud != nil && recipe != nil {
		am, err = ud.FindMetadata(recipe)

		if err == nil {
			buildPath = am.BuildPath
		}
	} else if len(target) == 0 {
		err = errors.New("Empty build path")
	}

	if err == nil {
		switch {
		case strings.HasSuffix(buildPath, ".apk"):
			flavor = data.BuildFlavorAndroid

		case strings.HasSuffix(buildPath, ".app"):
			flavor = data.BuildFlavorIos

		default:
			err = fmt.Errorf("File extension of build path at %q is not recognized", buildPath)
		}
	}

	if err == nil {
		return buildPath, flavor, nil
	}

	return "", "", err
}

func (ua *UploadAction) determineUploadToken(recipe *data.Recipe) (string, error) {
	uploadToken := ua.options.UploadToken

	if len(uploadToken) == 0 && recipe != nil {
		uploadToken = recipe.UploadToken
	}

	if len(uploadToken) == 0 {
		uploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	err := data.ValidateUploadToken(uploadToken)

	if err == nil {
		return uploadToken, nil
	}

	return "", err
}

func (ua *UploadAction) displaySummary(uploadToken, buildPath, buildPayloadPath string) {
	ua.ioStreams.Printf("\n")

	ua.summarize("Build path:", buildPath)
	ua.summarize("Git branch:", ua.options.GitBranch)
	ua.summarize("Git commit:", ua.options.GitCommit)
	ua.summarizeSecure("Upload token:", uploadToken)
	ua.summarize("Variant name:", ua.options.VariantName)

	if ua.options.Verbose {
		ua.ioStreams.Printf("\n")

		ua.summarize("Build payload path:", buildPayloadPath)
		ua.summarize("CI git branch:", ua.ciInfo.GitBranch)
		ua.summarize("CI git commit:", ua.ciInfo.GitCommit)
		ua.summarize("CI provider:", ua.ciInfo.Provider.String())
		ua.summarize("Git access:", ua.gitInfo.Access.String())
		ua.summarize("Inferred git branch:", ua.gitInfo.Branch)
		ua.summarize("Inferred git commit:", ua.gitInfo.Commit)
	}

	ua.ioStreams.Printf("\n")
}

func (ua *UploadAction) errorContentType() string {
	return lib.JsonContentType
}

func (ua *UploadAction) makeBuildPayloadPath(workingPath, buildPath string, flavor data.BuildFlavor) string {
	buildName := filepath.Base(buildPath)

	switch flavor {
	case data.BuildFlavorIos:
		return filepath.Join(workingPath, buildName+".zip")

	default:
		return buildPath
	}
}

func (ua *UploadAction) makeBuildURL(flavor data.BuildFlavor) string {
	var (
		wrapperName    string
		wrapperVersion string
	)

	buildURL := ua.apiBuildEndpoint

	if len(buildURL) == 0 {
		buildURL = data.DefaultAPIBuildEndpoint
	}

	query := make(url.Values)

	if len(ua.wrapperName) > 0 || len(ua.wrapperVersion) > 0 {
		wrapperName = ua.wrapperName
		wrapperVersion = ua.wrapperVersion
	} else {
		wrapperName = data.AgentName
		wrapperVersion = data.AgentVersion
	}

	lib.AddIfNotEmpty(&query, "agentName", data.AgentNameOld) // for now…
	lib.AddIfNotEmpty(&query, "agentVersion", data.AgentVersion)
	lib.AddIfNotEmpty(&query, "arch", ua.runtimeInfo.Arch)
	lib.AddIfNotEmpty(&query, "ci", ua.ciInfo.Provider.String())
	lib.AddIfNotEmpty(&query, "ciGitBranch", ua.ciInfo.GitBranch)
	lib.AddIfNotEmpty(&query, "ciGitCommit", ua.ciInfo.GitCommit)
	lib.AddIfNotEmpty(&query, "flavor", string(flavor))
	lib.AddIfNotEmpty(&query, "gitAccess", ua.gitInfo.Access.String())
	lib.AddIfNotEmpty(&query, "gitBranch", ua.gitInfo.Branch)
	lib.AddIfNotEmpty(&query, "gitCommit", ua.gitInfo.Commit)
	lib.AddIfNotEmpty(&query, "platform", ua.runtimeInfo.Platform)
	lib.AddIfNotEmpty(&query, "userGitBranch", ua.options.GitBranch)
	lib.AddIfNotEmpty(&query, "userGitCommit", ua.options.GitCommit)
	lib.AddIfNotEmpty(&query, "variantName", ua.options.VariantName)
	lib.AddIfNotEmpty(&query, "wrapperName", wrapperName)
	lib.AddIfNotEmpty(&query, "wrapperVersion", wrapperVersion)

	buildURL += "?" + query.Encode()

	return buildURL
}

func (ua *UploadAction) makeErrorPayload(err error) string {
	var (
		payload        string
		wrapperName    string
		wrapperVersion string
	)

	if len(ua.wrapperName) > 0 || len(ua.wrapperVersion) > 0 {
		wrapperName = ua.wrapperName
		wrapperVersion = ua.wrapperVersion
	} else {
		wrapperName = data.AgentName
		wrapperVersion = data.AgentVersion
	}

	lib.AppendIfNotEmpty(&payload, "agentName", data.AgentNameOld) // for now…
	lib.AppendIfNotEmpty(&payload, "agentVersion", data.AgentVersion)
	lib.AppendIfNotEmpty(&payload, "arch", ua.runtimeInfo.Arch)
	lib.AppendIfNotEmpty(&payload, "ci", ua.ciInfo.Provider.String())
	lib.AppendIfNotEmpty(&payload, "ciGitBranch", ua.ciInfo.GitBranch)
	lib.AppendIfNotEmpty(&payload, "ciGitCommit", ua.ciInfo.GitCommit)
	lib.AppendIfNotEmpty(&payload, "message", err.Error())
	lib.AppendIfNotEmpty(&payload, "platform", ua.runtimeInfo.Platform)
	lib.AppendIfNotEmpty(&payload, "wrapperName", wrapperName)
	lib.AppendIfNotEmpty(&payload, "wrapperVersion", wrapperVersion)

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

func (ua *UploadAction) prepareWorkingPath() (string, error) {
	path := filepath.Join(os.TempDir(), fmt.Sprintf("waldo-upload-%d", os.Getpid()))

	err := os.RemoveAll(path)

	if err == nil {
		err = os.MkdirAll(path, 0755)
	}

	if err == nil {
		return path, nil
	}

	return "", err
}

func (ua *UploadAction) summarize(label, value string) {
	if len(value) > 0 {
		ua.ioStreams.Printf("%-20.20s %q\n", label, value)
	} else {
		ua.ioStreams.Printf("%-20.20s (none)\n", label)
	}
}

func (ua *UploadAction) summarizeSecure(label, value string) {
	if len(value) == 0 {
		ua.ioStreams.Printf("%-20.20s (none)\n", label)
	} else if !ua.options.Verbose {
		prefixLen := len(value)

		if prefixLen > 6 {
			prefixLen = 6
		}

		prefix := value[0:prefixLen]
		suffixLen := len(value) - len(prefix)
		secure := strings.Repeat("*", 32)

		value = prefix + secure[0:suffixLen]

		ua.ioStreams.Printf("%-20.20s %q\n", label, value)
	} else {
		ua.ioStreams.Printf("%-20.20s %q\n", label, value)
	}
}

func (ua *UploadAction) uploadBuild(uploadToken, buildPayloadPath string, flavor data.BuildFlavor, retryAllowed bool) (bool, error) {
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

	req.Header.Add("Authorization", ua.authorization(uploadToken))

	if contentType := ua.buildContentType(flavor); len(contentType) > 0 {
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

func (ua *UploadAction) uploadBuildWithRetry(uploadToken, buildPayloadPath string, flavor data.BuildFlavor) error {
	for attempts := 1; attempts <= data.MaxPostAttempts; attempts++ {
		retry, err := ua.uploadBuild(uploadToken, buildPayloadPath, flavor, attempts < data.MaxPostAttempts)

		if !retry || err == nil {
			return err
		}

		ua.ioStreams.EmitError(data.AgentPrefix, err)

		ua.ioStreams.Printf("\nFailed upload attempts: %d -- retrying…\n\n", attempts)
	}

	return nil
}

func (ua *UploadAction) uploadError(err error, uploadToken string, flavor data.BuildFlavor, retryAllowed bool) (bool, error) {
	url := ua.makeErrorURL()
	body := ua.makeErrorPayload(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return false, fmt.Errorf("Unable to upload error to Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", ua.authorization(uploadToken))
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

func (ua *UploadAction) uploadErrorWithRetry(err error, uploadToken string, flavor data.BuildFlavor) error {
	for attempts := 1; attempts <= data.MaxPostAttempts; attempts++ {
		retry, tmpErr := ua.uploadError(err, uploadToken, flavor, attempts < data.MaxPostAttempts)

		if !retry || tmpErr == nil {
			return tmpErr
		}

		ua.ioStreams.EmitError(data.AgentPrefix, err)

		ua.ioStreams.Printf("\nFailed upload error attempts: %d -- retrying…\n\n", attempts)
	}

	return nil
}

func (ua *UploadAction) userAgent(flavor data.BuildFlavor) string {
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
