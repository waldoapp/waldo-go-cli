package reactnative

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

type BuildInfo struct { // package.json
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Modes           []string
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath string, platform lib.Platform, ios *lib.IOStreams) (*BuildInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "package.json"))

	if err != nil {
		return nil, err
	}

	var bi BuildInfo

	if err = json.Unmarshal(data, &bi); err != nil {
		return nil, err
	}

	if platform != lib.PlatformUnknown && ios != nil {
		bi.Modes, err = detectModes(basePath, bi.Name, platform)
	}

	return &bi, nil
}

//-----------------------------------------------------------------------------

func detectModes(basePath, name string, platform lib.Platform) ([]string, error) {
	switch platform {
	case lib.PlatformAndroid:
		return detectModesForAndroid(basePath)

	case lib.PlatformIos:
		return detectModesForIos(basePath, name)

	default:
		return nil, fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func detectModesForAndroid(basePath string) ([]string, error) {
	androidPath := filepath.Join(basePath, "android")

	bi, err := gradle.DetectBuildInfo(androidPath, "app")

	if err != nil {
		return []string{"release"}, nil
	}

	modes := lib.CompactMap(bi.Variants, func(mode string) (string, bool) {
		return mode, strings.ToLower(mode) == "release"
	})

	return modes, nil
}

func detectModesForIos(basePath, name string) ([]string, error) {
	iosPath := filepath.Join(basePath, "ios")

	bi, err := xcode.DetectBuildInfo(iosPath, name+".xcodeproj")

	if err != nil {
		return []string{"Debug"}, nil
	}

	modes := lib.CompactMap(bi.Configurations, func(mode string) (string, bool) {
		return mode, strings.ToLower(mode) == "debug"
	})

	return modes, nil
}
