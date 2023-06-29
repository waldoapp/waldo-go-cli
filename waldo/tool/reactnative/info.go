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

type ReactNativeInfo struct {
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Modes           []string
}

//-----------------------------------------------------------------------------

func DetectReactNativeInfo(basePath string, platform lib.Platform, ios *lib.IOStreams) (*ReactNativeInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "package.json"))

	if err != nil {
		return nil, err
	}

	var rni ReactNativeInfo

	if err = json.Unmarshal(data, &rni); err != nil {
		return nil, err
	}

	if platform != lib.PlatformUnknown && ios != nil {
		rni.Modes, err = detectModes(basePath, rni.Name, platform)
	}

	return &rni, nil
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

	gi, err := gradle.DetectGradleInfo(androidPath, "app")

	if err != nil {
		return []string{"release"}, nil
	}

	modes := lib.CompactMap(gi.Variants, func(mode string) bool {
		return strings.ToLower(mode) == "release"
	})

	return modes, nil
}

func detectModesForIos(basePath, name string) ([]string, error) {
	iosPath := filepath.Join(basePath, "ios")

	xi, err := xcode.DetectXcodeInfo(iosPath, name+".xcodeproj")

	if err != nil {
		return []string{"Debug"}, nil
	}

	modes := lib.CompactMap(xi.Configurations(), func(mode string) bool {
		return strings.ToLower(mode) == "debug"
	})

	return modes, nil
}
