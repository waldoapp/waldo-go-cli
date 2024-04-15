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

type Builder struct {
	Flavor string `yaml:"flavor,omitempty"`
}

type BuildInfo struct {
	//
	// From pubspec.yaml:
	//
	Name string `yaml:"name"`
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath string, platform lib.Platform) (*BuildInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "pubspec.yaml"))

	if err != nil {
		return nil, err
	}

	bi := &BuildInfo{}

	if err := tpw.DecodeFromYAML(data, bi); err != nil {
		return nil, err
	}

	return bi, nil
}

func IsPossibleContainer(path string) (bool, bool) {
	pubspecPath := filepath.Join(path, "pubspec.yaml")

	if !lib.IsRegularFile(pubspecPath) {
		return false, false
	}

	androidDirPath := filepath.Join(path, "android")
	iosDirPath := filepath.Join(path, "ios")

	hasAndroidProject, _ := gradle.IsPossibleContainer(androidDirPath)
	_, hasIosProject := xcode.IsPossibleContainer(iosDirPath)

	return hasAndroidProject, hasIosProject
}

func MakeBuilder(basePath string, verbose bool, ios *lib.IOStreams) (*Builder, string, lib.Platform, error) {
	platform, err := DeterminePlatform(verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	bi, err := DetectBuildInfo(basePath, platform)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	b := &Builder{}

	switch platform {
	case lib.PlatformAndroid:
		err = b.configureGradle(bi, basePath, verbose, ios)

	case lib.PlatformIos:
		err = b.configureXcode(bi, basePath, verbose, ios)

	default:
		err = fmt.Errorf("Unknown build platform: %q", platform)
	}

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	return b, bi.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.formatTarget(platform)

	ios.Printf("\nDetermining build path for %v\n", target)

	buildPath, err := b.determineBuildPath(basePath, platform)

	if err != nil {
		return "", err
	}

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	if clean {
		ios.Printf("\nCleaning %v\n", target)

		ios.Println(dashes)

		if err = b.clean(basePath, verbose, ios); err != nil {
			return "", err
		}

		ios.Println(dashes)
	}

	ios.Printf("\nBuilding %v\n", target)

	ios.Println(dashes)

	if err = b.build(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %v\n", target)

	return b.verifyBuildPath(buildPath, platform)
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

func (b *Builder) configureGradle(bi *BuildInfo, basePath string, verbose bool, ios *lib.IOStreams) error {
	androidPath := filepath.Join(basePath, "android")

	gbi, err := gradle.DetectBuildInfo(androidPath, "app")

	if err != nil {
		return err
	}

	flavor, err := DetermineFlavor(gbi.Variants, verbose, ios)

	if err != nil {
		return err
	}

	b.Flavor = flavor

	return nil
}

func (b *Builder) configureXcode(bi *BuildInfo, basePath string, verbose bool, ios *lib.IOStreams) error {
	iosPath := filepath.Join(basePath, "ios")

	xbi, err := xcode.DetectBuildInfo(iosPath, "Runner.xcodeproj")

	if err != nil {
		return err
	}

	flavors := lib.CompactMap(xbi.Configurations, func(flavor string) (string, bool) {
		switch strings.ToLower(flavor) {
		case "profile", "release":
			return "", false

		default:
			return flavor, true
		}
	})

	flavor, err := DetermineFlavor(flavors, verbose, ios)

	if err != nil {
		return err
	}

	b.Flavor = flavor

	return nil
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

func (b *Builder) formatTarget(platform lib.Platform) string {
	result := fmt.Sprintf("Flutter (%v)", platform)

	lib.AppendIfNotEmpty(&result, "flavor", b.Flavor, ": ", ", ")

	return result
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
		return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
	}

	return path, nil
}
