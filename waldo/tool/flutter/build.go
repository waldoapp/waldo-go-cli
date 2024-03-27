package flutter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type Builder struct {
	Flavor string `yaml:"flavor,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleContainer(path string) bool {
	pubspecPath := filepath.Join(path, "pubspec.yaml")
	androidDirPath := filepath.Join(path, "android")
	iosDirPath := filepath.Join(path, "ios")

	if !lib.IsRegularFile(pubspecPath) {
		return false
	}

	hasAndroidProject := gradle.IsPossibleContainer(androidDirPath)
	hasIosProject := xcode.IsPossibleContainer(iosDirPath)

	return hasAndroidProject || hasIosProject
}

func MakeBuilder(basePath string, verbose bool, ios *lib.IOStreams) (*Builder, string, lib.Platform, error) {
	platform, err := determinePlatform(verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all supported build flavors\n")

	fi, err := DetectBuildInfo(basePath, platform, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	flavor, err := determineFlavor(fi.Flavors, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	fb := &Builder{Flavor: flavor}

	return fb, fi.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.FormatTarget(platform)

	ios.Printf("\nDetermining build path for %s\n", target)

	buildPath, err := b.determineBuildPath(basePath, platform)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if clean {
		if err = b.clean(basePath, verbose, ios); err != nil {
			return "", err
		}
	}

	if err = b.build(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s\n", target)

	return b.verifyBuildPath(buildPath, platform)
}

func (b *Builder) FormatTarget(platform lib.Platform) string {
	result := string(platform)

	if len(b.Flavor) > 0 {
		result += " (" + b.Flavor + ")"
	}

	return result
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "flavor", b.Flavor, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func (b *Builder) androidBuildArgs() []string {
	args := []string{"apk"}

	switch strings.ToLower(b.Flavor) {
	case "debug":
		args = append(args, "--debug")

	case "profile":
		args = append(args, "--profile")

	case "release":
		args = append(args, "--release")

	default:
		args = append(args, "--flavor", b.Flavor)
	}

	return args
}

func (b *Builder) build(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"build"}

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, b.androidBuildArgs()...)

	case lib.PlatformIos:
		args = append(args, b.iosBuildArgs()...)

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

func (b *Builder) clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	args := []string{"clean"}

	if verbose {
		args = append(args, "--verbose")
	}

	task := lib.NewTask("flutter", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) determineBuildPath(basePath string, platform lib.Platform) (string, error) {
	relPath := ""

	switch platform {
	case lib.PlatformAndroid:
		relPath = "build/app/outputs/flutter-apk/app-" + strings.ToLower(b.Flavor) + ".apk"

	case lib.PlatformIos:
		relPath = "build/ios/iphonesimulator/Runner.app"

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}

	return filepath.Join(basePath, relPath), nil
}

func (b *Builder) iosBuildArgs() []string {
	return []string{"ios", "--simulator", "--no-codesign"}
}

func (b *Builder) verifyBuildPath(path string, platform lib.Platform) (string, error) {
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
