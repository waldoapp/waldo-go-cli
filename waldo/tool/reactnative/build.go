package reactnative

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type Builder struct {
	Mode string `yaml:"mode,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleContainer(path string) bool {
	packagePath := filepath.Join(path, "package.json")
	androidDirPath := filepath.Join(path, "android")
	iosDirPath := filepath.Join(path, "ios")

	if !lib.IsRegularFile(packagePath) {
		return false
	}

	bi, err := DetectBuildInfo(path, lib.PlatformUnknown, nil)

	if err != nil {
		return false
	}

	_, exFound := bi.Dependencies["expo"]
	_, rnFound := bi.Dependencies["react-native"]

	if exFound || !rnFound {
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

	ios.Printf("\nFinding all supported build modes\n")

	bi, err := DetectBuildInfo(basePath, platform, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	mode, err := determineMode(bi.Modes, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	b := &Builder{Mode: mode}

	return b, bi.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.FormatTarget(platform)

	bi, err := DetectBuildInfo(basePath, platform, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %s\n", target)

	buildPath, err := b.determineBuildPath(basePath, bi.Name, platform, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if clean {
		if err = b.clean(basePath, bi.Name, platform, verbose, ios); err != nil {
			return "", err
		}
	}

	if err = b.build(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s\n", target)

	return b.verifyBuildPath(buildPath, bi.Name, platform, ios)
}

func (b *Builder) FormatTarget(platform lib.Platform) string {
	result := string(platform)

	if len(b.Mode) > 0 {
		result += " (" + b.Mode + ")"
	}

	return result
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "mode", b.Mode, "=", ", ")

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

func (b *Builder) clean(basePath, name string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.Builder{
			Module:  "app",
			Variant: b.Mode}

		return gb.Clean(filepath.Join(basePath, "android"), verbose, ios)

	case lib.PlatformIos:
		xb := &xcode.Builder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: b.Mode}

		return xb.Clean(filepath.Join(basePath, "ios"), verbose, ios)

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (b *Builder) determineBuildPath(basePath, name string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.Builder{
			Module:  "app",
			Variant: b.Mode}

		return gb.DetermineBuildPath(filepath.Join(basePath, "android"), ios)

	case lib.PlatformIos:
		xb := &xcode.Builder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: b.Mode}

		return xb.DetermineBuildPath(filepath.Join(basePath, "ios"), ios)

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}
}

func (b *Builder) iosBuildArgs() []string {
	return []string{"build-ios", "--mode", b.Mode}
}

func (b *Builder) verifyBuildPath(path, name string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	switch platform {
	case lib.PlatformAndroid:
		gb := &gradle.Builder{
			Module:  "app",
			Variant: b.Mode}

		return gb.VerifyBuildPath(path, ios)

	case lib.PlatformIos:
		xb := &xcode.Builder{
			Workspace:     name + ".xcworkspace",
			Scheme:        name,
			Configuration: b.Mode}

		return xb.VerifyBuildPath(path, ios)

	default:
		return "", fmt.Errorf("Unknown build platform: %q", platform)
	}

	return path, nil
}
