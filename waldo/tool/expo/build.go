package expo

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
	Variant       string `yaml:"variant,omitempty"`
	Scheme        string `yaml:"scheme,omitempty"`
	Configuration string `yaml:"configuration,omitempty"`
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

	if !exFound || !rnFound {
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
		Variant: b.Variant}
}

func (b *Builder) XcodeBuilder() *xcode.Builder {
	return &xcode.Builder{
		Workspace:     b.Scheme + ".xcworkspace",
		Scheme:        b.Scheme,
		Configuration: b.Configuration}
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

	ios.Printf("\nPrebuilding %v\n", target)

	ios.Println(dashes)

	if err = b.prebuild(basePath, platform, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nBuilding %v\n", target)

	ios.Println(dashes)

	if err = b.build(basePath, platform, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %v\n", target)

	return b.verifyBuildPath(buildPath, platform, ios)
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "variant", b.Variant, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "scheme", b.Scheme, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "configuration", b.Configuration, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func (b *Builder) build(basePath string, platform lib.Platform, clean, verbose bool, ios *lib.IOStreams) error {
	switch platform {
	case lib.PlatformAndroid:
		_, err := b.GradleBuilder().Build(filepath.Join(basePath, "android"), clean, verbose, ios)

		return err

	case lib.PlatformIos:
		_, err := b.XcodeBuilder().Build(filepath.Join(basePath, "ios"), clean, verbose, ios)

		return err

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

	variants := lib.CompactMap(gbi.Variants, func(variant string) (string, bool) {
		return variant, strings.ToLower(variant) != "debug"
	})

	variant, err := gradle.DetermineVariant("app", variants, verbose, ios)

	if err != nil {
		return err
	}

	b.Variant = variant

	return nil
}

func (b *Builder) configureXcode(bi *BuildInfo, basePath string, verbose bool, ios *lib.IOStreams) error {
	iosPath := filepath.Join(basePath, "ios")

	baseName := b.matchBaseName(iosPath, bi.Name)

	if len(baseName) == 0 {
		return fmt.Errorf("No base name found matching %q", bi.Name)
	}

	xbi, err := xcode.DetectBuildInfo(iosPath, baseName+".xcodeproj")

	if err != nil {
		return err
	}

	schemes := lib.CompactMap(xbi.Schemes, func(scheme string) (string, bool) {
		return scheme, scheme == baseName
	})

	configs := lib.CompactMap(xbi.Configurations, func(config string) (string, bool) {
		return config, strings.ToLower(config) != "debug"
	})

	workspace := baseName + ".xcworkspace"

	scheme, err := xcode.DetermineScheme(workspace, schemes, verbose, ios)

	if err != nil {
		return err
	}

	config, err := xcode.DetermineConfiguration(workspace, configs, verbose, ios)

	if err != nil {
		return err
	}

	b.Scheme = scheme
	b.Configuration = config

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
	result := fmt.Sprintf("Expo (%v)", platform)

	lib.AppendIfNotEmpty(&result, "variant", b.Variant, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "scheme", b.Scheme, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "configuration", b.Configuration, ": ", ", ")

	return result
}

func (b *Builder) matchBaseName(basePath, name string) string {
	wsSuffix := ".xcworkspace"
	wsPaths := lib.FindDirectoryPathsMatching(filepath.Join(basePath, "*"+wsSuffix))

	baseName := ""

	for _, wsPath := range wsPaths {
		wsName, found := strings.CutSuffix(filepath.Base(wsPath), wsSuffix)

		if !found || strings.ToLower(wsName) != name {
			continue
		}

		baseName = wsName
		break
	}

	if len(baseName) == 0 {
		return ""
	}

	prSuffix := ".xcodeproj"
	prPaths := lib.FindDirectoryPathsMatching(filepath.Join(basePath, "*"+prSuffix))

	for _, prPath := range prPaths {
		prName, found := strings.CutSuffix(filepath.Base(prPath), prSuffix)

		if !found || prName != baseName {
			continue
		}

		return baseName
	}

	return ""
}

func (b *Builder) prebuild(basePath string, platform lib.Platform, verbose bool, ios *lib.IOStreams) error {
	args := []string{"expo", "prebuild", "--platform"}

	switch platform {
	case lib.PlatformAndroid:
		args = append(args, "android")

	case lib.PlatformIos:
		args = append(args, "ios")

	default:
		return fmt.Errorf("Unknown build platform: %q", platform)
	}

	args = append(args, "--no-install")

	task := lib.NewTask("npx", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
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
