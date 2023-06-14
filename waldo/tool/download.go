package tool

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

const (
	agentAssetBaseURL = "https://github.com/waldoapp/waldo-go-agent/releases"

	agentMaxGetAttempts = 2
)

//-----------------------------------------------------------------------------

type AgentDownloader struct {
	agentPath    string
	assetURL     string
	assetVersion string
	errorPrefix  string
	ioStreams    *lib.IOStreams
	runtimeInfo  *lib.RuntimeInfo
	verbose      bool
	workingPath  string
}

//-----------------------------------------------------------------------------

func NewAgentDownloader(assetVersion, errorPrefix string, verbose bool, ioStreams *lib.IOStreams, runtimeInfo *lib.RuntimeInfo) *AgentDownloader {
	return &AgentDownloader{
		assetVersion: assetVersion,
		errorPrefix:  errorPrefix,
		ioStreams:    ioStreams,
		runtimeInfo:  runtimeInfo,
		verbose:      verbose}
}

//-----------------------------------------------------------------------------

func (ad *AgentDownloader) Cleanup() {
	if len(ad.workingPath) > 0 {
		os.RemoveAll(ad.workingPath)
	}
}

func (ad *AgentDownloader) Download() (string, error) {
	if err := ad.prepareSource(); err != nil {
		return "", err
	}

	if err := ad.prepareTarget(); err != nil {
		return "", err
	}

	if err := ad.downloadAgentWithRetry(); err != nil {
		return "", err
	}

	return ad.agentPath, nil
}

//-----------------------------------------------------------------------------

func (ad *AgentDownloader) checkStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return fmt.Errorf("Unable to download Waldo Agent, HTTP status: %d", status)
	}

	return nil
}

func (ad *AgentDownloader) determineAgentPath() string {
	agentName := "waldo-agent"

	if ad.runtimeInfo.Platform == "Windows" {
		agentName += ".exe"
	}

	return filepath.Join(ad.workingPath, "waldo-agent")
}

func (ad *AgentDownloader) determineAssetURL() string {
	platform := strings.ToLower(string(ad.runtimeInfo.Platform))
	arch := strings.ToLower(string(ad.runtimeInfo.Arch))

	assetName := fmt.Sprintf("waldo-agent-%s-%s", platform, arch)

	if platform == "windows" {
		assetName += ".exe"
	}

	assetBaseURL := agentAssetBaseURL

	if ad.assetVersion != "latest" {
		assetBaseURL += "/download/" + ad.assetVersion
	} else {
		assetBaseURL += "/latest/download"
	}

	return assetBaseURL + "/" + assetName
}

func (ad *AgentDownloader) determineWorkingPath() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("WaldoGoCLI-%d", os.Getpid()))
}

func (ad *AgentDownloader) downloadAgent(retryAllowed bool) (bool, error) {
	ad.ioStreams.Printf("\nDownloading Waldo Agent…\n\n")

	client := &http.Client{}

	req, err := http.NewRequest("GET", ad.assetURL, nil)

	if err != nil {
		return false, fmt.Errorf("Unable to download Waldo Agent, error: %v, url: %s", err, ad.assetURL)
	}

	if ad.verbose {
		lib.DumpRequest(ad.ioStreams, req, false)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return retryAllowed, fmt.Errorf("Unable to download Waldo Agent, error: %v, url: %s", err, ad.assetURL)
	}

	if ad.verbose {
		lib.DumpResponse(ad.ioStreams, rsp, false)
	}

	defer rsp.Body.Close()

	err = ad.checkStatus(rsp)

	if err != nil {
		return retryAllowed && lib.ShouldRetry(rsp), err
	}

	err = ad.saveResponseBody(rsp, ad.agentPath)

	if err != nil {
		return false, fmt.Errorf("Unable to download Waldo Agent, error: %v, url: %s", err, ad.assetURL)
	}

	return false, nil
}

func (ad *AgentDownloader) downloadAgentWithRetry() error {
	for attempts := 1; attempts <= agentMaxGetAttempts; attempts++ {
		retry, err := ad.downloadAgent(attempts < agentMaxGetAttempts)

		if !retry || err == nil {
			return err
		}

		ad.ioStreams.EmitError(ad.errorPrefix, err)

		ad.ioStreams.Printf("\nFailed download attempts: %d -- retrying…\n", attempts)
	}

	return nil
}

func (ad *AgentDownloader) prepareSource() error {
	ad.assetURL = ad.determineAssetURL()

	return nil
}

func (ad *AgentDownloader) prepareTarget() error {
	ad.workingPath = ad.determineWorkingPath()
	ad.agentPath = ad.determineAgentPath()

	if err := os.RemoveAll(ad.workingPath); err != nil {
		return err
	}

	return os.MkdirAll(ad.workingPath, 0755)
}

func (ad *AgentDownloader) saveResponseBody(rsp *http.Response, path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0775)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, rsp.Body)

	return err
}
