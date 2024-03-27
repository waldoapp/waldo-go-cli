package expo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type BuildInfo struct {
	//
	// From package.json:
	//
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	//
	// Extra info:
	//
	GradleVariant      string
	XcodeScheme        string
	XcodeConfiguration string
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath string, platform lib.Platform, ios *lib.IOStreams) (*BuildInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "package.json"))

	if err != nil {
		return nil, err
	}

	bi := &BuildInfo{}

	if err = json.Unmarshal(data, bi); err != nil {
		return nil, err
	}

	if ios != nil {
		switch platform {
		case lib.PlatformAndroid:
			err = bi.detectExtraForAndroid(basePath, ios)

		case lib.PlatformIos:
			err = bi.detectExtraForIos(basePath, bi.Name, ios)

		default:
			err = fmt.Errorf("Unknown build platform: %q", platform)
		}

		if err != nil {
			return nil, err
		}
	}

	return bi, nil
}

//-----------------------------------------------------------------------------

func (bi *BuildInfo) detectExtraForAndroid(basePath string, ios *lib.IOStreams) error {
	androidPath := filepath.Join(basePath, "android")

	gbi, err := gradle.DetectBuildInfo(androidPath, "app")

	if err != nil {
		return err
	}

	variants := lib.CompactMap(gbi.Variants, func(variant string) (string, bool) {
		return variant, strings.ToLower(variant) != "debug"
	})

	ios.Printf("*** Detected variants: %s", variants)

	bi.GradleVariant = variants[0]

	return nil
}

func (bi *BuildInfo) detectExtraForIos(basePath, name string, ios *lib.IOStreams) error {
	iosPath := filepath.Join(basePath, "ios")

	xbi, err := xcode.DetectBuildInfo(iosPath, name+".xcodeproj")

	if err != nil {
		return err
	}

	schemes := lib.CompactMap(xbi.Schemes, func(scheme string) (string, bool) {
		return scheme, strings.ToLower(scheme) == name
	})

	configs := lib.CompactMap(xbi.Configurations, func(config string) (string, bool) {
		return config, strings.ToLower(config) != "debug"
	})

	ios.Printf("*** Detected schemes: %s, configurations: %s", schemes, configs)

	bi.XcodeConfiguration = configs[0]
	bi.XcodeScheme = schemes[0]

	return nil
}
