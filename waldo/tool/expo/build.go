package expo

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/gradle"
	"github.com/waldoapp/waldo-go-cli/waldo/tool/xcode"
)

type Builder struct {
	Gradle *gradle.Builder `yaml:"gradle,omitempty"`
	Xcode  *xcode.Builder  `yaml:"xcode,omitempty"`
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

	if !exFound || !rnFound {
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

	bi, err := DetectBuildInfo(basePath, platform, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	b := &Builder{}

	switch platform {
	case lib.PlatformAndroid:
		b.Gradle = &gradle.Builder{
			Module:  "app",
			Variant: bi.GradleVariant}

	case lib.PlatformIos:
		b.Xcode = &xcode.Builder{
			Workspace:     bi.Name + ".xcworkspace",
			Project:       "",
			Scheme:        bi.XcodeScheme,
			Configuration: bi.XcodeConfiguration}

	default:
		return nil, "", platform, fmt.Errorf("Unknown build platform: %q", platform)
	}

	return b, bi.Name, platform, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := b.FormatTarget()

	bi, err := DetectBuildInfo(basePath, platform, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %s\n", target)

	buildPath, err := b.determineBuildPath(basePath, bi.Name, platform, ios)

	if err != nil {
		return "", err
	}

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Printf("\nPrebuilding %s\n", target)

	ios.Println(dashes)

	if err = b.prebuild(basePath, bi.Name, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nBuilding %s\n", target)

	ios.Println(dashes)

	if err = b.build(basePath, bi.Name, platform, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s\n", target)

	return b.verifyBuildPath(buildPath, ios)
}

func (b *Builder) FormatTarget() string {
	if b.Gradle != nil {
		return b.Gradle.FormatTarget()
	}

	if b.Xcode != nil {
		return b.Xcode.FormatTarget()
	}

	return ""
}

func (b *Builder) Summarize() string {
	if b.Gradle != nil {
		return b.Gradle.Summarize()
	}

	if b.Xcode != nil {
		return b.Xcode.Summarize()
	}

	return ""
}

//-----------------------------------------------------------------------------

func (b *Builder) androidBuildArgs() []string {
	return []string{"run:android", "--variant", "release"}
}

func (b *Builder) androidPrebuildArgs() []string {
	return []string{"--platform", "android"}
}

func (b *Builder) build(basePath, name string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) error {
	if platform == lib.PlatformIos {
		return b.buildForIos(basePath, name, clean, verbose, ios)
	}

	args := []string{"expo"}

	args = append(args, b.androidBuildArgs()...)

	if clean {
		args = append(args, "--no-build-cache")
	}

	args = append(args, "--no-bundler", "--no-install")

	task := lib.NewTask("npx", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) buildForIos(basePath, name string, clean, verbose bool, ios *lib.IOStreams) error {
	_, err := b.Xcode.Build(filepath.Join(basePath, "ios"), clean, verbose, ios)

	return err
}

func (b *Builder) determineBuildPath(basePath, name string, platform lib.Platform, ios *lib.IOStreams) (string, error) {
	if b.Gradle != nil {
		return b.Gradle.DetermineBuildPath(filepath.Join(basePath, "android"), ios)
	}

	if b.Xcode != nil {
		return b.Xcode.DetermineBuildPath(filepath.Join(basePath, "ios"), ios)
	}

	return "", fmt.Errorf("Unknown build platform: %q", platform)
}

func (b *Builder) iosBuildArgs() []string {
	return []string{"run:ios", "--configuration", "Release"}
}

func (b *Builder) iosPrebuildArgs() []string {
	return []string{"--platform", "ios"}
}

func (b *Builder) prebuild(basePath, name string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"expo", "prebuild"}

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, b.androidPrebuildArgs()...)

	case lib.PlatformIos:
		args = append(args, b.iosPrebuildArgs()...)

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}

	args = append(args, "--no-install")

	task := lib.NewTask("npx", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) verifyBuildPath(path string, ios *lib.IOStreams) (string, error) {
	if b.Gradle != nil {
		return b.Gradle.VerifyBuildPath(path, ios)
	}

	if b.Xcode != nil {
		return b.Xcode.VerifyBuildPath(path, ios)
	}

	return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
}
