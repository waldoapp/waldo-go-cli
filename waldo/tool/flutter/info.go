package flutter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/lib/tpw"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type BuildInfo struct { // pubspec.yaml
	Name    string `yaml:"name"`
	Flavors []string
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath string, platform lib.Platform, ios *lib.IOStreams) (*BuildInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "pubspec.yaml"))

	if err != nil {
		return nil, err
	}

	var bi BuildInfo

	if err := tpw.DecodeFromYAML(data, &bi); err != nil {
		return nil, err
	}

	bi.Flavors, err = detectFlavors(basePath, platform, ios)

	return &bi, nil
}

//-----------------------------------------------------------------------------

func detectFlavors(basePath string, platform lib.Platform, ios *lib.IOStreams) ([]string, error) {
	switch platform {
	case lib.PlatformAndroid:
		return detectFlavorsForAndroid(basePath, ios)

	case lib.PlatformIos:
		return detectFlavorsForIos(basePath)

	default:
		return nil, fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func detectFlavorsForAndroid(basePath string, ios *lib.IOStreams) ([]string, error) {
	androidPath := filepath.Join(basePath, "android")

	bi, err := gradle.DetectBuildInfo(androidPath, "app")

	if err != nil {
		return []string{"debug", "release"}, nil
	}

	return bi.Variants, nil
}

func detectFlavorsForIos(basePath string) ([]string, error) {
	iosPath := filepath.Join(basePath, "ios")

	bi, err := xcode.DetectBuildInfo(iosPath, "Runner.xcodeproj")

	if err != nil {
		return []string{"Debug"}, nil
	}

	flavors := lib.CompactMap(bi.Configurations, func(flavor string) (string, bool) {
		switch strings.ToLower(flavor) {
		case "profile", "release":
			return "", false

		default:
			return flavor, true
		}
	})

	return flavors, nil
}
