package reactnative

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type ReactNativeBuilder struct {
	Mode string `yaml:"mode,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleReactNativeContainer(path string) bool {
	packagePath := filepath.Join(path, "package.json")
	androidDirPath := filepath.Join(path, "android")
	iosDirPath := filepath.Join(path, "ios")

	if !lib.IsRegularFile(packagePath) {
		return false
	}

	pkg, err := DetectReactNativeInfo(path, lib.PlatformUnknown, nil)

	if err != nil {
		return false
	}

	_, found := pkg.Dependencies["react-native"]

	if !found {
		return false
	}

	hasAndroidProject := gradle.IsPossibleGradleContainer(androidDirPath)
	hasIosProject := xcode.IsPossibleXcodeContainer(iosDirPath)

	return hasAndroidProject || hasIosProject
}

func MakeReactNativeBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*ReactNativeBuilder, string, lib.Platform, error) {
	platform, err := determinePlatform(verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all supported build modes…\n")

	rni, err := DetectReactNativeInfo(absPath, platform, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	mode, err := determineMode(rni.Modes, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	rnb := &ReactNativeBuilder{Mode: mode}

	return rnb, rni.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (rnb *ReactNativeBuilder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := rnb.formatTarget(platform)

	rni, err := DetectReactNativeInfo(basePath, platform, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %s…\n", target)

	buildPath, err := rnb.determineBuildPath(basePath, rni.Name, platform, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s…\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if clean {
		if err = rnb.clean(basePath, rni.Name, platform, verbose, ios); err != nil {
			return "", err
		}
	}

	if err = rnb.build(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s…\n", target)

	return rnb.verifyBuildPath(buildPath, rni.Name, platform, ios)
}

func (rnb *ReactNativeBuilder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "mode", rnb.Mode, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func (rnb *ReactNativeBuilder) androidBuildArgs() []string {
	return []string{"build-android", "--mode", rnb.Mode, "--no-packager"}
}

func (rnb *ReactNativeBuilder) build(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"--yes", "react-native"}
	env := lib.CurrentEnvironment()

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, rnb.androidBuildArgs()...)

	case lib.PlatformIos:
		args = append(args, rnb.iosBuildArgs()...)

		env["FORCE_BUNDLING"] = "1" // this is the magic!

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}

	if verbose {
		args = append(args, "--verbose")
	}

	task := lib.NewTask("npx", args...)

	task.Cwd = basePath
	task.Env = env
	task.IOStreams = ios

	return task.Execute()
}

func (rnb *ReactNativeBuilder) clean(basePath, name string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.GradleBuilder{
			Module:  "app",
			Variant: rnb.Mode}

		return gb.Clean(filepath.Join(basePath, "android"), verbose, ios)

	case lib.PlatformIos:
		xb := &xcode.XcodeBuilder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: rnb.Mode}

		return xb.Clean(filepath.Join(basePath, "ios"), verbose, ios)

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (rnb *ReactNativeBuilder) determineBuildPath(basePath, name string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.GradleBuilder{
			Module:  "app",
			Variant: rnb.Mode}

		return gb.DetermineBuildPath(filepath.Join(basePath, "android"), ios)

	case lib.PlatformIos:
		xb := &xcode.XcodeBuilder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: rnb.Mode}

		return xb.DetermineBuildPath(filepath.Join(basePath, "ios"), ios)

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (rnb *ReactNativeBuilder) formatTarget(platform lib.Platform) string {
	result := string(platform)

	if len(rnb.Mode) > 0 {
		result += " (" + rnb.Mode + ")"
	}

	return result
}

func (rnb *ReactNativeBuilder) iosBuildArgs() []string {
	return []string{"build-ios", "--mode", rnb.Mode}
}

func (rnb *ReactNativeBuilder) verifyBuildPath(path, name string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.GradleBuilder{
			Module:  "app",
			Variant: rnb.Mode}

		return gb.VerifyBuildPath(path, ios)

	case lib.PlatformIos:
		xb := &xcode.XcodeBuilder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: rnb.Mode}

		return xb.VerifyBuildPath(path, ios)

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}

	return path, nil
}

//-----------------------------------------------------------------------------

func reactNativePath() string {
	path, _, err := lib.NewTask("command", "-v", "react-native").Run()

	if err == nil && len(path) > 0 {
		return path
	}

	path, _, err = lib.NewTask("npx", "--yes", "which", "react-native").Run()

	if err == nil && len(path) > 0 {
		return path
	}

	return "react-native"
}
