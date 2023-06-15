package flutter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type FlutterBuilder struct {
	Flavor string `yaml:"flavor,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleFlutterContainer(path string) bool {
	pubspecPath := filepath.Join(path, "pubspec.yaml")
	androidDirPath := filepath.Join(path, "android")
	iosDirPath := filepath.Join(path, "ios")

	if !lib.IsRegularFile(pubspecPath) {
		return false
	}

	hasAndroidProject := gradle.IsPossibleGradleContainer(androidDirPath)
	hasIosProject := xcode.IsPossibleXcodeContainer(iosDirPath)

	return hasAndroidProject || hasIosProject
}

func MakeFlutterBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*FlutterBuilder, string, lib.Platform, error) {
	platform, err := determinePlatform(verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all supported build flavors…\n")

	fi, err := DetectFlutterInfo(absPath, platform, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	flavor, err := determineFlavor(fi.Flavors, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	fb := &FlutterBuilder{
		Flavor: flavor}

	return fb, fi.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (fb *FlutterBuilder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := fb.formatTarget(platform)

	ios.Printf("\nDetermining build path for %s…\n", target)

	buildPath, err := fb.determineBuildPath(basePath, platform)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s…\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if clean {
		if err = fb.clean(basePath, verbose, ios); err != nil {
			return "", err
		}
	}

	if err = fb.build(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s…\n", target)

	return fb.verifyBuildPath(buildPath, platform)
}

func (fb *FlutterBuilder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "flavor", fb.Flavor, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func (fb *FlutterBuilder) androidBuildArgs() []string {
	args := []string{"apk"}

	switch strings.ToLower(fb.Flavor) {
	case "debug":
		args = append(args, "--debug")

	case "profile":
		args = append(args, "--profile")

	case "release":
		args = append(args, "--release")

	default:
		args = append(args, "--flavor", fb.Flavor)
	}

	return args
}

func (fb *FlutterBuilder) build(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"build"}

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, fb.androidBuildArgs()...)

	case lib.PlatformIos:
		args = append(args, fb.iosBuildArgs()...)

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}

	args = append(args, "--no-tree-shake-icons")

	if verbose {
		args = append(args, "--verbose")
	}

	task := lib.NewTask("flutter", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (fb *FlutterBuilder) clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	args := []string{"clean"}

	if verbose {
		args = append(args, "--verbose")
	}

	task := lib.NewTask("flutter", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (fb *FlutterBuilder) determineBuildPath(basePath string, platform lib.Platform) (string, error) {
	relPath := ""

	switch platform {
	case lib.PlatformAndroid:
		relPath = "build/app/outputs/flutter-apk/app-" + strings.ToLower(fb.Flavor) + ".apk"

	case lib.PlatformIos:
		relPath = "build/ios/iphonesimulator/Runner.app"

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}

	return filepath.Join(basePath, relPath), nil
}

func (fb *FlutterBuilder) formatTarget(platform lib.Platform) string {
	result := string(platform)

	if len(fb.Flavor) > 0 {
		result += " (" + fb.Flavor + ")"
	}

	return result
}

func (fb *FlutterBuilder) iosBuildArgs() []string {
	return []string{"ios", "--simulator", "--no-codesign"}
}

func (fb *FlutterBuilder) verifyBuildPath(path string, platform lib.Platform) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		if !lib.IsRegularFile(path) {
			return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
		}

	case lib.PlatformIos:
		if !lib.IsDirectory(path) {
			return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
		}

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}

	return path, nil
}
