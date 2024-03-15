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

type FlutterInfo struct { // pubspec.yaml
	Name    string `yaml:"name"`
	Flavors []string
}

//-----------------------------------------------------------------------------

func DetectFlutterInfo(basePath string, platform lib.Platform, ios *lib.IOStreams) (*FlutterInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "pubspec.yaml"))

	if err != nil {
		return nil, err
	}

	var fi FlutterInfo

	if err := tpw.DecodeFromYAML(data, &fi); err != nil {
		return nil, err
	}

	fi.Flavors, err = detectFlavors(basePath, platform, ios)

	return &fi, nil
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

	gi, err := gradle.DetectGradleInfo(androidPath, "app")

	if err != nil {
		return []string{"debug", "release"}, nil
	}

	return gi.Variants, nil
}

func detectFlavorsForIos(basePath string) ([]string, error) {
	iosPath := filepath.Join(basePath, "ios")

	xi, err := xcode.DetectXcodeInfo(iosPath, "Runner.xcodeproj")

	if err != nil {
		return []string{"Debug"}, nil
	}

	flavors := lib.CompactMap(xi.Configurations(), func(flavor string) (string, bool) {
		switch strings.ToLower(flavor) {
		case "profile", "release":
			return "", false

		default:
			return flavor, true
		}
	})

	return flavors, nil
}
