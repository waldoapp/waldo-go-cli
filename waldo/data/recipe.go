package data

import (
	"github.com/waldoapp/waldo-go-cli/waldo/tool"
)

type BuildFlavor string

const (
	BuildFlavorAndroid = "Android"
	BuildFlavorIos     = "iOS"
	BuildFlavorUnknown = "Unknown"
)

type Recipe struct {
	Name               string                   `yaml:"recipe"`
	AppName            string                   `yaml:"app_name"`
	Flavor             BuildFlavor              `yaml:"platform"`
	UploadToken        string                   `yaml:"upload_token"`
	BasePath           string                   `yaml:"build_root"` // relative to Configuration.BasePath
	BeforeBuild        *ShellScript             `yaml:"before_build,omitempty"`
	CustomBuilder      *tool.CustomBuilder      `yaml:"custom_build,omitempty"`
	ExpoBuilder        *tool.ExpoBuilder        `yaml:"expo_build,omitempty"`
	FlutterBuilder     *tool.FlutterBuilder     `yaml:"flutter_build,omitempty"`
	GradleBuilder      *tool.GradleBuilder      `yaml:"gradle_build,omitempty"`
	ReactNativeBuilder *tool.ReactNativeBuilder `yaml:"reactnative_build,omitempty"`
	XcodeBuilder       *tool.XcodeBuilder       `yaml:"xcode_build,omitempty"`
	AfterBuild         *ShellScript             `yaml:"after_build,omitempty"`
}

type ShellScript struct {
	Script      string            `yaml:"script"`
	Environment map[string]string `yaml:"environment,omitempty"`
}

//-----------------------------------------------------------------------------

func (r *Recipe) BuildTool() tool.BuildTool {
	if r.CustomBuilder != nil {
		return tool.BuildToolCustom
	}

	if r.ExpoBuilder != nil {
		return tool.BuildToolExpo
	}

	if r.FlutterBuilder != nil {
		return tool.BuildToolFlutter
	}

	if r.GradleBuilder != nil {
		return tool.BuildToolGradle
	}

	if r.ReactNativeBuilder != nil {
		return tool.BuildToolReactNative
	}

	if r.XcodeBuilder != nil {
		return tool.BuildToolXcode
	}

	return tool.BuildToolCustom
}

func (r *Recipe) Summarize() string {
	switch r.BuildTool() {
	case tool.BuildToolCustom:
		return r.CustomBuilder.Summarize()

	case tool.BuildToolExpo:
		return r.ExpoBuilder.Summarize()

	case tool.BuildToolFlutter:
		return r.FlutterBuilder.Summarize()

	case tool.BuildToolGradle:
		return r.GradleBuilder.Summarize()

	case tool.BuildToolReactNative:
		return r.ReactNativeBuilder.Summarize()

	case tool.BuildToolXcode:
		return r.XcodeBuilder.Summarize()

	default:
		return ""
	}
}
