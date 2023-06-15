package xcode

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type XcodeBuilder struct {
	Workspace     string `yaml:"workspace,omitempty"`
	Project       string `yaml:"project,omitempty"`
	Scheme        string `yaml:"scheme,omitempty"`
	Configuration string `yaml:"configuration,omitempty"`
}

//-----------------------------------------------------------------------------

func IsPossibleXcodeContainer(path string) bool {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	return lib.HasDirectoryMatching(workPath) || lib.HasDirectoryMatching(projPath)
}

func MakeXcodeBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*XcodeBuilder, string, lib.Platform, error) {
	ios.Printf("\nSearching for Xcode workspaces and projects in %q…\n", relPath)

	fileName, err := determineProject(findXcodeProjects(absPath), verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all Xcode schemes and configurations in %q…\n", fileName)

	xi, err := DetectXcodeInfo(absPath, fileName)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	scheme, err := determineScheme(fileName, xi.Schemes(), verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	configuration, err := determineConfiguration(fileName, xi.Configurations(), verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	var (
		project   string
		workspace string
	)

	if strings.HasSuffix(fileName, ".xcworkspace") {
		workspace = fileName
	} else {
		project = fileName
	}

	xb := &XcodeBuilder{
		Workspace:     workspace,
		Project:       project,
		Scheme:        scheme,
		Configuration: configuration}

	return xb, xi.Name(), lib.PlatformIos, nil
}

//-----------------------------------------------------------------------------

func (xb *XcodeBuilder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	target := xb.formatTarget()

	ios.Printf("\nDetecting build settings for %s…\n", target)

	settings, err := xb.detectBuildSettings(basePath, ios)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %s…\n", target)

	buildPath, err := xb.determineBuildPath(settings)

	if err != nil {
		return "", err
	}

	ios.Printf("\nBuilding %s…\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if err = xb.build(basePath, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %s…\n", target)

	return xb.verifyBuildPath(buildPath)
}

func (xb *XcodeBuilder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "workspace", xb.Workspace, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "project", xb.Project, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "scheme", xb.Scheme, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "configuration", xb.Configuration, "=", ", ")

	return summary
}

//-----------------------------------------------------------------------------

func findXcodeProjects(path string) []string {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	workspaces := lib.FindDirectoryPathsMatching(workPath)
	projects := lib.FindDirectoryPathsMatching(projPath)

	return append(workspaces, projects...)
}

//-----------------------------------------------------------------------------

func (xb *XcodeBuilder) build(basePath string, clean, verbose bool, ios *lib.IOStreams) error {
	args := xb.commonBuildArgs()

	if !verbose {
		args = append(args, "-quiet")
	}

	if clean {
		args = append(args, "clean")
	}

	args = append(args, "build")

	task := lib.NewTask("xcodebuild", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (xb *XcodeBuilder) commonBuildArgs() []string {
	args := []string{}

	if len(xb.Workspace) > 0 {
		args = append(args, "-workspace", xb.Workspace)
	} else {
		args = append(args, "-project", xb.Project)
	}

	if len(xb.Scheme) > 0 {
		args = append(args, "-scheme", xb.Scheme)
	}

	if len(xb.Configuration) > 0 {
		args = append(args, "-configuration", xb.Configuration)
	}

	return append(args, "-sdk", "iphonesimulator")
}

func (xb *XcodeBuilder) detectBuildSettings(basePath string, ios *lib.IOStreams) (map[string]string, error) {
	args := xb.commonBuildArgs()

	args = append(args, "build")

	args = append([]string{"-showBuildSettings"}, args...)

	task := lib.NewTask("xcodebuild", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	results, _, err := task.Run()

	if err != nil {
		return nil, err
	}

	return xb.parseBuildSettings(results), nil
}

func (xb *XcodeBuilder) determineBuildPath(settings map[string]string) (string, error) {
	buildName := settings["FULL_PRODUCT_NAME"]
	buildDir := settings["TARGET_BUILD_DIR"]

	if len(buildDir) == 0 || len(buildName) == 0 {
		return "", errors.New("Unable to determine build path")
	}

	return filepath.Join(buildDir, buildName), nil
}

func (xb *XcodeBuilder) formatTarget() string {
	result := ""

	if len(xb.Workspace) > 0 {
		result += xb.Workspace
	} else {
		result += xb.Project
	}

	lib.AppendIfNotEmpty(&result, "scheme", xb.Scheme, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "configuration", xb.Configuration, ": ", ", ")

	return result
}

func (xb *XcodeBuilder) parseBuildSettings(text string) map[string]string {
	settings := make(map[string]string)

	for _, line := range strings.Split(text, "\n") {
		pair := strings.SplitN(line, "=", 2)

		if len(pair) != 2 {
			continue
		}

		key := strings.TrimSpace(pair[0])
		value := strings.TrimSpace(pair[1])

		if len(key) > 0 && len(value) > 0 {
			settings[key] = value
		}
	}

	return settings
}

func (xb *XcodeBuilder) verifyBuildPath(path string) (string, error) {
	if !lib.IsDirectory(path) {
		return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
	}

	return path, nil
}
