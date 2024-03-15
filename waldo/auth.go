package waldo

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
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

	aa.ioStreams.Printf("\nAuthenticating with user token %q…\n", aa.options.UserToken)

	fullName, err := api.AuthenticateUser(aa.options.UserToken, aa.options.Verbose, aa.ioStreams)

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
