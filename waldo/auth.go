package waldo

import (
	"fmt"
	"net/http"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type AuthOptions struct {
	UserToken string
	Verbose   bool
}

type AuthAction struct {
	ioStreams   *lib.IOStreams
	options     *AuthOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewAuthAction(options *AuthOptions, ioStreams *lib.IOStreams) *AuthAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &AuthAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (aa *AuthAction) Perform() error {
	if err := data.ValidateUserToken(aa.options.UserToken); err != nil {
		return err
	}

	fullName, err := aa.authenticateUser()

	if err != nil {
		return err
	}

	profile, _, err := data.SetupProfile(data.CreateKindIfNeeded)

	if err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	profile.UserToken = aa.options.UserToken

	profile.MarkDirty()

	if err := profile.Save(); err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	aa.ioStreams.Printf("\nUser %q successfully authenticated -- credentials saved to %q\n", fullName, profile.Path())

	return nil
}

//-----------------------------------------------------------------------------

func (aa *AuthAction) authenticateUser() (string, error) {
	aa.ioStreams.Printf("\nAuthenticating with user token %q…\n", aa.options.UserToken)

	client := &http.Client{}

	req, err := http.NewRequest("GET", data.CoreAPIUserEndpoint, nil)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", aa.options.UserToken))
	req.Header.Add("User-Agent", data.FullVersion())

	if aa.options.Verbose {
		lib.DumpRequest(aa.ioStreams, req, true)
	}

	rsp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	defer rsp.Body.Close()

	if aa.options.Verbose {
		lib.DumpResponse(aa.ioStreams, rsp, true)
	}

	status := rsp.StatusCode

	if status < 200 || status > 299 {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", rsp.Status)
	}

	ar, err := data.ParseAuthResponse(rsp)

	if err != nil {
		return "", fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	return ar.FullName(), nil
}
