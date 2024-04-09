package data

import (
	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/expo"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/flutter"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/ionic"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/reactnative"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type Recipe struct {
	Name               string               `yaml:"recipe"`
	AppName            string               `yaml:"app_name"`
	AppID              string               `yaml:"app_id"`
	Platform           lib.Platform         `yaml:"platform"`
	BasePath           string               `yaml:"build_root"` // relative to Configuration.BasePath
	ExpoBuilder        *expo.Builder        `yaml:"expo_build,omitempty"`
	FlutterBuilder     *flutter.Builder     `yaml:"flutter_build,omitempty"`
	GradleBuilder      *gradle.Builder      `yaml:"gradle_build,omitempty"`
	IonicBuilder       *ionic.Builder       `yaml:"ionic_build,omitempty"`
	ReactNativeBuilder *reactnative.Builder `yaml:"reactnative_build,omitempty"`
	XcodeBuilder       *xcode.Builder       `yaml:"xcode_build,omitempty"`
}

//-----------------------------------------------------------------------------

func (r *Recipe) BuildTool() tool.BuildTool {
	if r.ExpoBuilder != nil {
		return tool.BuildToolExpo
	}

	if r.FlutterBuilder != nil {
		return tool.BuildToolFlutter
	}

	if r.GradleBuilder != nil {
		return tool.BuildToolGradle
	}

	if r.IonicBuilder != nil {
		return tool.BuildToolIonic
	}

	if r.ReactNativeBuilder != nil {
		return tool.BuildToolReactNative
	}

	if r.XcodeBuilder != nil {
		return tool.BuildToolXcode
	}

	return tool.BuildToolUnknown
}

func (r *Recipe) Summarize() string {
	switch r.BuildTool() {
	case tool.BuildToolExpo:
		return r.ExpoBuilder.Summarize()

	case tool.BuildToolFlutter:
		return r.FlutterBuilder.Summarize()

	case tool.BuildToolGradle:
		return r.GradleBuilder.Summarize()

	case tool.BuildToolIonic:
		return r.IonicBuilder.Summarize()

	case tool.BuildToolReactNative:
		return r.ReactNativeBuilder.Summarize()

	case tool.BuildToolXcode:
		return r.XcodeBuilder.Summarize()

	default:
		return ""
	}
}
