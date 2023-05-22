package waldo

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type TriggerOptions struct {
	GitCommit   string
	RuleName    string
	UploadToken string
	Verbose     bool
}

type TriggerAction struct {
	apiTriggerEndpoint string
	ciInfo             *lib.CIInfo
	ioStreams          *lib.IOStreams
	options            *TriggerOptions
	runtimeInfo        *lib.RuntimeInfo
	wrapperName        string
	wrapperVersion     string
}

//-----------------------------------------------------------------------------

func NewTriggerAction(options *TriggerOptions, ioStreams *lib.IOStreams, overrides map[string]string) *TriggerAction {
	ciInfo := lib.DetectCIInfo(false)
	runtimeInfo := lib.DetectRuntimeInfo()

	return &TriggerAction{
		apiTriggerEndpoint: overrides["apiTriggerEndpoint"],
		ciInfo:             ciInfo,
		ioStreams:          ioStreams,
		options:            options,
		runtimeInfo:        runtimeInfo,
		wrapperName:        overrides["wrapperName"],
		wrapperVersion:     overrides["wrapperVersion"]}
}

//-----------------------------------------------------------------------------

func (ta *TriggerAction) Perform() error {
	uploadToken, err := ta.determineUploadToken()

	if err == nil {
		err = ta.triggerRunWithRetry(uploadToken)
	}

	if err == nil {
		ta.ioStreams.Printf("\nRun successfully triggered on Waldo!\n")
	}

	return err
}

//-----------------------------------------------------------------------------

func (ta *TriggerAction) authorization(uploadToken string) string {
	return fmt.Sprintf("Upload-Token %s", uploadToken)
}

func (ta *TriggerAction) checkTriggerStatus(rsp *http.Response) error {
	status := rsp.StatusCode

	if status == 401 {
		return errors.New("Upload token is invalid or missing!")
	}

	if status < 200 || status > 299 {
		return fmt.Errorf("Unable to trigger run on Waldo, HTTP status: %d", status)
	}

	return nil
}

func (ta *TriggerAction) contentType() string {
	return lib.JsonContentType
}

func (ta *TriggerAction) determineUploadToken() (string, error) {
	uploadToken := ta.options.UploadToken

	if len(uploadToken) == 0 {
		uploadToken = os.Getenv("WALDO_UPLOAD_TOKEN")
	}

	err := data.ValidateUploadToken(uploadToken)

	if err == nil {
		return uploadToken, nil
	}

	return "", err
}

func (ta *TriggerAction) makePayload() string {
	var (
		payload        string
		wrapperName    string
		wrapperVersion string
	)

	if len(ta.wrapperName) > 0 || len(ta.wrapperVersion) > 0 {
		wrapperName = ta.wrapperName
		wrapperVersion = ta.wrapperVersion
	} else {
		wrapperName = data.AgentName
		wrapperVersion = data.AgentVersion
	}

	lib.AppendIfNotEmpty(&payload, "agentName", data.AgentNameOld) // for now…
	lib.AppendIfNotEmpty(&payload, "agentVersion", data.AgentVersion)
	lib.AppendIfNotEmpty(&payload, "arch", ta.runtimeInfo.Arch)
	lib.AppendIfNotEmpty(&payload, "ci", ta.ciInfo.Provider.String())
	lib.AppendIfNotEmpty(&payload, "gitSha", ta.options.GitCommit)
	lib.AppendIfNotEmpty(&payload, "platform", ta.runtimeInfo.Platform)
	lib.AppendIfNotEmpty(&payload, "ruleName", ta.options.RuleName)
	lib.AppendIfNotEmpty(&payload, "wrapperName", wrapperName)
	lib.AppendIfNotEmpty(&payload, "wrapperVersion", wrapperVersion)

	payload = "{" + payload + "}"

	return payload
}

func (ta *TriggerAction) makeURL() string {
	triggerURL := ta.apiTriggerEndpoint

	if len(triggerURL) == 0 {
		triggerURL = data.DefaultAPITriggerEndpoint
	}

	return triggerURL
}

func (ta *TriggerAction) triggerRun(uploadToken string, retryAllowed bool) (bool, error) {
	ta.ioStreams.Printf("\nTriggering run on Waldo…\n")

	url := ta.makeURL()
	body := ta.makePayload()

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(body))

	if err != nil {
		return false, fmt.Errorf("Unable to trigger run on Waldo, error: %v, url: %s", err, url)
	}

	req.Header.Add("Authorization", ta.authorization(uploadToken))
	req.Header.Add("Content-Type", ta.contentType())
	req.Header.Add("User-Agent", ta.userAgent())

	if ta.options.Verbose {
		lib.DumpRequest(ta.ioStreams, req, true)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return retryAllowed, fmt.Errorf("Unable to trigger run on Waldo, error: %v, url: %s", err, url)
	}

	if ta.options.Verbose {
		lib.DumpResponse(ta.ioStreams, rsp, true)
	}

	defer rsp.Body.Close()

	return retryAllowed && lib.ShouldRetry(rsp), ta.checkTriggerStatus(rsp)
}

func (ta *TriggerAction) triggerRunWithRetry(uploadToken string) error {
	for attempts := 1; attempts <= data.MaxPostAttempts; attempts++ {
		retry, err := ta.triggerRun(uploadToken, attempts < data.MaxPostAttempts)

		if !retry || err == nil {
			return err
		}

		ta.ioStreams.EmitError(data.AgentPrefix, err)

		ta.ioStreams.Printf("\nFailed trigger attempts: %d -- retrying…\n\n", attempts)
	}

	return nil

}

func (ta *TriggerAction) userAgent() string {
	ci := ta.ciInfo.Provider.String()

	if ci == "Unknown" {
		ci = "Go CLI"
	}

	version := ta.wrapperVersion

	if len(version) == 0 {
		version = data.AgentVersion
	}

	return fmt.Sprintf("Waldo %s v%s", ci, version)
}
