package waldo

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type AuthOptions struct {
	APIToken string
	Verbose  bool
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
	if err := data.ValidateAPIToken(aa.options.APIToken); err != nil {
		return err
	}

	aa.ioStreams.Printf("\nAuthenticating with API token %q\n", aa.options.APIToken)

	fullName, err := api.AuthenticateUser(aa.options.APIToken, aa.options.Verbose, aa.ioStreams)

	if err != nil {
		return err
	}

	profile, _, err := data.SetupProfile(data.CreateKindIfNeeded)

	if err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	profile.APIToken = aa.options.APIToken

	profile.MarkDirty()

	if err := profile.Save(); err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	aa.ioStreams.Printf("\nUser %q successfully authenticated -- credentials saved to %q\n", fullName, profile.Path())

	return nil
}
