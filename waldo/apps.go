package waldo

import (
	"fmt"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/api"
	"github.com/waldoapp/waldo-go-cli/waldo/data"
)

type AppsOptions struct {
	AndroidOnly bool
	IosOnly     bool
}

type AppsAction struct {
	ioStreams   *lib.IOStreams
	options     *AppsOptions
	runtimeInfo *lib.RuntimeInfo
}

//-----------------------------------------------------------------------------

func NewAppsAction(options *AppsOptions, ioStreams *lib.IOStreams) *AppsAction {
	runtimeInfo := lib.DetectRuntimeInfo()

	return &AppsAction{
		ioStreams:   ioStreams,
		options:     options,
		runtimeInfo: runtimeInfo}
}

//-----------------------------------------------------------------------------

func (aa *AppsAction) Perform() error {
	profile, _, err := data.SetupProfile(data.CreateKindNever)

	if err != nil {
		return fmt.Errorf("Unable to authenticate user, error: %v", err)
	}

	platform := aa.determinePlatform()

	items, err := api.FetchApps(profile.APIToken, platform, false, aa.ioStreams)

	if err != nil {
		return err
	}

	aa.ioStreams.Printf("%-24.24v  %-20.20v  %-8.8v\n", "APP NAME", "APP ID", "PLATFORM")

	for _, item := range items {
		appName := aa.formatString(item.AppName, "(none)")
		appID := aa.formatString(item.AppID, "(none)")
		platform := item.Platform

		aa.ioStreams.Printf("%-24.24v  %-20.20v  %-8.8v\n", appName, appID, platform)

	}

	return nil
}

//-----------------------------------------------------------------------------

func (aa *AppsAction) determinePlatform() lib.Platform {
	if aa.options.AndroidOnly && aa.options.IosOnly {
		return lib.PlatformUnknown
	}

	if aa.options.AndroidOnly {
		return lib.PlatformAndroid
	}

	if aa.options.IosOnly {
		return lib.PlatformIos
	}

	return lib.PlatformUnknown
}

func (aa *AppsAction) formatString(value, defaultValue string) string {
	if len(value) > 0 {
		return value
	}

	return defaultValue
}
