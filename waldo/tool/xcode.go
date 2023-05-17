package tool

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
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

func FindXcodeProjects(path string) []string {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	workspaces := lib.FindDirectoryPathsMatching(workPath)
	projects := lib.FindDirectoryPathsMatching(projPath)

	return append(workspaces, projects...)
}

func IsPossibleXcodeContainer(path string) bool {
	workPath := filepath.Join(path, "*.xcworkspace")
	projPath := filepath.Join(path, "*.xcodeproj")

	return lib.HasDirectoryMatching(workPath) || lib.HasDirectoryMatching(projPath)
}

func MakeXcodeBuilder(absPath, relPath string, verbose bool, ios *lib.IOStreams) (*XcodeBuilder, error) {
	ios.Printf("\nSearching for Xcode workspaces and projects in %q…\n", relPath)

	project, err := determineProject(FindXcodeProjects(absPath), verbose, ios)

	if err != nil {
		return nil, err
	}

	ios.Printf("\nFinding all Xcode schemes and configurations in %q…\n", project)

	xi, err := detectXcodeInfo(absPath, project)

	if err != nil {
		return nil, err
	}

	scheme, err := determineScheme(project, xi, verbose, ios)

	if err != nil {
		return nil, err
	}

	configuration, err := determineConfiguration(project, xi, verbose, ios)

	if err != nil {
		return nil, err
	}

	xb := NewXcodeBuilder(project, scheme, configuration)

	return xb, nil
}

func NewXcodeBuilder(fileName, scheme, configuration string) *XcodeBuilder {
	var (
		project   string
		workspace string
	)

	if strings.HasSuffix(fileName, ".xcworkspace") {
		workspace = fileName
	} else {
		project = fileName
	}

	return &XcodeBuilder{
		Workspace:     workspace,
		Project:       project,
		Scheme:        scheme,
		Configuration: configuration}
}

//-----------------------------------------------------------------------------

func (xb *XcodeBuilder) Build(basePath string, verbose bool, ios *lib.IOStreams) error {
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

	args = append(args, "-sdk", "iphonesimulator")

	if !verbose {
		args = append(args, "-quiet")
	}

	args = append(args, "build")

	return lib.ExecDir(basePath, "xcodebuild", args...)
}

func (xb *XcodeBuilder) Summarize() string {
	summary := ""

	if len(xb.Workspace) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "workspace: " + xb.Workspace
	}

	if len(xb.Project) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "project: " + xb.Project
	}

	if len(xb.Scheme) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "scheme: " + xb.Scheme
	}

	if len(xb.Configuration) > 0 {
		if len(summary) > 0 {
			summary += ", "
		}

		summary += "configuration: " + xb.Configuration
	}

	return summary
}

//-----------------------------------------------------------------------------

func askConfiguration(configurations []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(configurations)

	idx, _ := pr.ReadChoose(
		"Available Xcode configurations",
		configurations,
		"Choose a configuration",
		0,
		true)

	return configurations[idx]
}

func askProject(projects []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(projects)

	idx, _ := pr.ReadChoose(
		"Available Xcode workspaces and projects",
		projects,
		"Choose a workspace or project",
		0,
		true)

	return projects[idx]
}

func askScheme(schemes []string, ios *lib.IOStreams) string {
	pr := ios.PromptReader()

	sort.Strings(schemes)

	idx, _ := pr.ReadChoose(
		"Available Xcode schemes",
		schemes,
		"Choose a scheme",
		0,
		true)

	return schemes[idx]
}

func determineConfiguration(project string, xi *xcodeInfo, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(xi.configurations) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode configuration found in %q\n", project)
		}

		return askConfiguration(xi.configurations, ios), nil
	}

	if len(xi.configurations) == 1 {
		if verbose {
			ios.Printf("\nOnly one Xcode configuration found in %q: %q\n", project, xi.configurations[0])
		}

		return xi.configurations[0], nil
	}

	if verbose {
		ios.Printf("\nNo Xcode configurations found in %q\n", project)
	}

	return "", nil
}

func determineProject(projects []string, verbose bool, ios *lib.IOStreams) (string, error) {
	projects = lib.Map(projects, func(project string) string {
		return filepath.Base(project)
	})

	if len(projects) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode workspace or project found\n")
		}

		return askProject(projects, ios), nil
	}

	if len(projects) == 1 {
		if verbose {
			ios.Printf("\nOnly one Xcode workspace or project found: %q\n", projects[0])
		}

		return projects[0], nil
	}

	return "", errors.New("No Xcode workspaces or projects found")
}

func determineScheme(project string, xi *xcodeInfo, verbose bool, ios *lib.IOStreams) (string, error) {
	if len(xi.schemes) > 1 {
		if verbose {
			ios.Printf("\nMore than one Xcode scheme found in %q\n", project)
		}

		return askScheme(xi.schemes, ios), nil
	}

	if len(xi.schemes) == 1 {
		if verbose {
			ios.Printf("\nOnly one Xcode scheme found in %q: %q\n", project, xi.schemes[0])
		}

		return xi.schemes[0], nil
	}

	return "", fmt.Errorf("No Xcode schemes found in %q", project)
}

//-----------------------------------------------------------------------------

type xcodeInfo struct {
	name           string
	configurations []string
	schemes        []string
}

//-----------------------------------------------------------------------------

func detectXcodeInfo(basePath, fileName string) (*xcodeInfo, error) {
	xi := &xcodeInfo{}

	err := xi.populate(basePath, fileName)

	if err != nil {
		return nil, err
	}

	return xi, nil
}

//-----------------------------------------------------------------------------

func (xi *xcodeInfo) populate(basePath, fileName string) error {
	if strings.HasSuffix(fileName, ".xcworkspace") {
		return xi.populateFromWorkspace(basePath, fileName)
	}

	return xi.populateFromProject(basePath, fileName)
}

func (xi *xcodeInfo) populateFromProject(basePath, project string) error {
	jsonString, _, err := lib.RunDir(basePath, "xcodebuild", "-list", "-json", "-project", project)

	if err == nil {
		rawJson := lib.ParseTopLevelJsonObject([]byte(jsonString))

		if rawJson != nil {
			project := lib.ParseJsonObject(rawJson["project"])

			if project != nil {
				xi.configurations = lib.ParseJsonStringArray(project["configurations"])
				xi.name = lib.ParseJsonString(project["name"])
				xi.schemes = lib.ParseJsonStringArray(project["schemes"])

				return nil
			}
		}
	}

	return err
}

func (xi *xcodeInfo) populateFromWorkspace(basePath, workspace string) error {
	jsonString, _, err := lib.RunDir(basePath, "xcodebuild", "-list", "-json", "-workspace", workspace)

	if err == nil {
		rawJson := lib.ParseTopLevelJsonObject([]byte(jsonString))

		if rawJson != nil {
			workspace := lib.ParseJsonObject(rawJson["workspace"])

			if workspace != nil {
				xi.name = lib.ParseJsonString(workspace["name"])
				xi.schemes = lib.ParseJsonStringArray(workspace["schemes"])

				return nil
			}
		}
	}

	return err
}
