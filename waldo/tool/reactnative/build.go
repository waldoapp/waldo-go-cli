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

type Builder struct {
	Mode string `yaml:"mode,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type BuildInfo struct {
	//
	// From package.json:
	//
	Name            string            `json:"name"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath string, platform lib.Platform) (*BuildInfo, error) {
	data, err := os.ReadFile(filepath.Join(basePath, "package.json"))

	if err != nil {
		return nil, err
	}

	bi := &BuildInfo{}

	if err = json.Unmarshal(data, bi); err != nil {
		return nil, err
	}

	return bi, nil
}

func IsPossibleContainer(path string) (bool, bool) {
	packagePath := filepath.Join(path, "package.json")

	if !lib.IsRegularFile(packagePath) {
		return false, false
	}

	bi, err := DetectBuildInfo(path, lib.PlatformUnknown)

	if err != nil {
		return false, false
	}

	_, exFound := bi.Dependencies["expo"]
	_, rnFound := bi.Dependencies["react-native"]

	if exFound || !rnFound {
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

func (b *Builder) GradleBuilder() *gradle.Builder {
	return &gradle.Builder{
		Module:  "app",
		Variant: b.Mode}
}

func (b *Builder) XcodeBuilder() *xcode.Builder {
	return &xcode.Builder{
		Workspace:     b.Name + ".xcworkspace",
		Scheme:        b.Name,
		Configuration: b.Mode}
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.formatTarget(platform)

	ios.Printf("\nDetermining build path for %v\n", target)

	buildPath, err := b.determineBuildPath(basePath, platform, ios)

	if err != nil {
		return "", err
	}

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	if clean {
		ios.Printf("\nCleaning %v\n", target)

		ios.Println(dashes)

		if err = b.clean(basePath, platform, verbose, ios); err != nil {
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

	return b.verifyBuildPath(buildPath, platform, ios)
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "mode", b.Mode, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "name", b.Name, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func (b *Builder) androidBuildArgs() []string {
	return []string{"build-android", "--mode", b.Mode, "--no-packager"}
}

func (b *Builder) build(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"--no-install", "react-native"}
	env := lib.CurrentEnvironment()

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, b.androidBuildArgs()...)

	case lib.PlatformIos:
		args = append(args, b.iosBuildArgs()...)

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

func (b *Builder) clean(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	switch platform {
	case lib.PlatformAndroid:
		return b.GradleBuilder().Clean(filepath.Join(basePath, "android"), verbose, ios)

	case lib.PlatformIos:
		return b.XcodeBuilder().Clean(filepath.Join(basePath, "ios"), verbose, ios)

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (b *Builder) configureGradle(bi *BuildInfo, basePath string, verbose bool, ios *lib.IOStreams) error {
	androidPath := filepath.Join(basePath, "android")

	gbi, err := gradle.DetectBuildInfo(androidPath, "app")

	if err != nil {
		return err
	}

	modes := lib.CompactMap(gbi.Variants, func(mode string) (string, bool) {
		return mode, strings.ToLower(mode) == "release"
	})

	mode, err := DetermineMode(modes, verbose, ios)

	if err != nil {
		return err
	}

	b.Mode = mode

	return nil
}

func (b *Builder) configureXcode(bi *BuildInfo, basePath string, verbose bool, ios *lib.IOStreams) error {
	iosPath := filepath.Join(basePath, "ios")

	xbi, err := xcode.DetectBuildInfo(iosPath, bi.Name+".xcodeproj")

	if err != nil {
		return err
	}

	modes := lib.CompactMap(xbi.Configurations, func(mode string) (string, bool) {
		return mode, strings.ToLower(mode) == "debug"
	})

	mode, err := DetermineMode(modes, verbose, ios)

	if err != nil {
		return err
	}

	b.Mode = mode
	b.Name = bi.Name

	return nil
}

func (b *Builder) determineBuildPath(basePath string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		return b.GradleBuilder().DetermineBuildPath(filepath.Join(basePath, "android"), ios)

	case lib.PlatformIos:
		return b.XcodeBuilder().DetermineBuildPath(filepath.Join(basePath, "ios"), ios)

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (b *Builder) formatTarget(platform lib.Platform) string {
	result := fmt.Sprintf("React Native (%v)", platform)

	lib.AppendIfNotEmpty(&result, "mode", b.Mode, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "name", b.Name, ": ", ", ")

	return result
}

func (b *Builder) iosBuildArgs() []string {
	return []string{"build-ios", "--mode", b.Mode}
}

func (b *Builder) verifyBuildPath(path string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		return b.GradleBuilder().VerifyBuildPath(path, ios)

	case lib.PlatformIos:
		return b.XcodeBuilder().VerifyBuildPath(path, ios)

	default:
		return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
	}
}
