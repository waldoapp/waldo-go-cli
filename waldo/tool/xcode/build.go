package xcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type Builder struct {
	Workspace     string `yaml:"workspace,omitempty"`
	Project       string `yaml:"project,omitempty"`
	Scheme        string `yaml:"scheme,omitempty"`
	Configuration string `yaml:"configuration,omitempty"`
}

type BuildInfo struct {
	Name           string
	Schemes        []string
	Configurations []string
}

type BuildSettings map[string]string

//-----------------------------------------------------------------------------

func DetectBuildInfo(basePath, fileName string) (*BuildInfo, error) {
	bi := &BuildInfo{}

	if err := bi.detect(basePath, fileName); err != nil {
		return nil, err
	}

	return bi, nil
}

func IsPossibleContainer(path string) bool {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	return lib.HasDirectoryMatching(workPath) || lib.HasDirectoryMatching(projPath)
}

func MakeBuilder(basePath string, verbose bool, ios *lib.IOStreams) (*Builder, string, lib.Platform, error) {
	ios.Printf("\nFinding all Xcode workspaces and projects in %q\n", lib.MakeRelativeToCWD(basePath))

	fileName, err := DetermineProject(findProjects(basePath), verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	ios.Printf("\nFinding all Xcode schemes and configurations in %q\n", fileName)

	bi, err := DetectBuildInfo(basePath, fileName)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	scheme, err := DetermineScheme(fileName, bi.Schemes, verbose, ios)

	if err != nil {
		return nil, "", lib.PlatformUnknown, err
	}

	configuration, err := DetermineConfiguration(fileName, bi.Configurations, verbose, ios)

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

	b := &Builder{
		Workspace:     workspace,
		Project:       project,
		Scheme:        scheme,
		Configuration: configuration}

	return b, bi.Name, lib.PlatformIos, nil
}

//-----------------------------------------------------------------------------

func (b *Builder) Build(basePath string, clean, verbose bool, ios *lib.IOStreams) (string, error) {
	buildPath, err := b.DetermineBuildPath(basePath, ios)

	if err != nil {
		return "", err
	}

	target := b.formatTarget()

	ios.Printf("\nBuilding %v\n", target)

	dashes := "\n" + strings.Repeat("-", 79) + "\n"

	ios.Println(dashes)

	if err = b.build(basePath, clean, verbose, ios); err != nil {
		return "", err
	}

	ios.Println(dashes)

	ios.Printf("\nVerifying build path for %v\n", target)

	return b.verifyBuildPath(buildPath)
}

func (b *Builder) Clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	ios.Printf("\nCleaning %v\n", b.formatTarget())

	return b.clean(basePath, verbose, ios)
}

func (b *Builder) DetermineBuildPath(basePath string, ios *lib.IOStreams) (string, error) {
	target := b.formatTarget()

	ios.Printf("\nDetecting build settings for %v\n", target)

	settings, err := b.detectBuildSettings(basePath)

	if err != nil {
		return "", err
	}

	ios.Printf("\nDetermining build path for %v\n", target)

	return b.determineBuildPath(settings)
}

func (b *Builder) Summarize() string {
	summary := ""

	lib.AppendIfNotEmpty(&summary, "workspace", b.Workspace, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "project", b.Project, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "scheme", b.Scheme, "=", ", ")
	lib.AppendIfNotEmpty(&summary, "configuration", b.Configuration, "=", ", ")

	return summary
}

func (b *Builder) VerifyBuildPath(basePath string, ios *lib.IOStreams) (string, error) {
	ios.Printf("\nVerifying build path for %v\n", b.formatTarget())

	return b.verifyBuildPath(basePath)
}

//-----------------------------------------------------------------------------

type SettingsResult struct {
	Action        string        `json:"action"`
	BuildSettings BuildSettings `json:"buildSettings"`
	Target        string        `json:"target"`
}

//-----------------------------------------------------------------------------

func findProjects(path string) []string {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	workspaces := lib.FindDirectoryPathsMatching(workPath)
	projects := lib.FindDirectoryPathsMatching(projPath)

	return append(workspaces, projects...)
}

//-----------------------------------------------------------------------------

func (b *Builder) build(basePath string, clean, verbose bool, ios *lib.IOStreams) error {
	args := b.commonArgs()

	args = append(args, "-sdk", "iphonesimulator")

	if verbose {
		args = append(args, "-verbose")
	} else {
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

func (b *Builder) clean(basePath string, verbose bool, ios *lib.IOStreams) error {
	args := b.commonArgs()

	if verbose {
		args = append(args, "-verbose")
	} else {
		args = append(args, "-quiet")
	}

	args = append(args, "clean")

	task := lib.NewTask("xcodebuild", args...)

	task.Cwd = basePath
	task.IOStreams = ios

	return task.Execute()
}

func (b *Builder) commonArgs() []string {
	args := []string{}

	if len(b.Workspace) > 0 {
		args = append(args, "-workspace", b.Workspace)
	} else if len(b.Project) > 0 {
		args = append(args, "-project", b.Project)
	}

	if len(b.Scheme) > 0 {
		args = append(args, "-scheme", b.Scheme)
	}

	if len(b.Configuration) > 0 {
		args = append(args, "-configuration", b.Configuration)
	}

	return args
}

func (b *Builder) detectBuildSettings(basePath string) (BuildSettings, error) {
	args := []string{"-showBuildSettings"}

	args = append(args, b.commonArgs()...)

	args = append(args, "-sdk", "iphonesimulator")

	args = append(args, "-json", "build")

	task := lib.NewTask("xcodebuild", args...)

	task.Cwd = basePath

	data, _, err := task.RunRaw()

	var results []SettingsResult

	err = json.Unmarshal(data, &results)

	if err != nil {
		return nil, err
	}

	return results[0].BuildSettings, nil
}

func (b *Builder) determineBuildPath(settings BuildSettings) (string, error) {
	buildDir := settings["CONFIGURATION_BUILD_DIR"]
	buildName := settings["FULL_PRODUCT_NAME"]

	if len(buildDir) == 0 || len(buildName) == 0 {
		return "", errors.New("Unable to determine build path")
	}

	return filepath.Join(buildDir, buildName), nil
}

func (b *Builder) formatTarget() string {
	result := "Xcode"

	lib.AppendIfNotEmpty(&result, "workspace", b.Workspace, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "project", b.Project, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "scheme", b.Scheme, ": ", ", ")
	lib.AppendIfNotEmpty(&result, "configuration", b.Configuration, ": ", ", ")

	return result
}

func (b *Builder) verifyBuildPath(path string) (string, error) {
	if !lib.IsDirectory(path) {
		return "", fmt.Errorf("Unable to locate build path, expected: %q", path)
	}

	return path, nil
}

//-----------------------------------------------------------------------------

type ProjectResult struct {
	Project Project `json:"project"`
}

type Project struct {
	Name           string   `json:"name"`
	Schemes        []string `json:"schemes"`
	Configurations []string `json:"configurations"`
}

type WorkspaceResult struct {
	Workspace Workspace `json:"workspace"`
}

type Workspace struct {
	Name    string   `json:"name"`
	Schemes []string `json:"schemes"`
}

//-----------------------------------------------------------------------------

func (bi *BuildInfo) detect(basePath, fileName string) error {
	isWorkspace := strings.HasSuffix(fileName, ".xcworkspace")

	args := []string{"-list"}

	if isWorkspace {
		args = append(args, "-workspace")
	} else {
		args = append(args, "-project")
	}

	args = append(args, fileName, "-json")

	task := lib.NewTask("xcodebuild", args...)

	task.Cwd = basePath

	data, _, err := task.RunRaw()

	if err != nil {
		return err
	}

	if isWorkspace {
		var result WorkspaceResult

		err = json.Unmarshal(data, &result)

		if err = json.Unmarshal(data, &result); err != nil {
			return err
		}

		bi.Name = result.Workspace.Name
		bi.Schemes = result.Workspace.Schemes
	} else {
		var result ProjectResult

		if err = json.Unmarshal(data, &result); err != nil {
			return err
		}

		bi.Name = result.Project.Name
		bi.Schemes = result.Project.Schemes
		bi.Configurations = result.Project.Configurations
	}

	return nil
}
