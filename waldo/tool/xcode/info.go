package xcode

import (
	"encoding/json"
	"strings"

	"github.com/waldoapp/waldo-go-cli/lib"
)

type XcodeInfo struct {
	Project   *ProjectInfo   `json:"project"`
	Workspace *WorkspaceInfo `json:"workspace"`
}

type ProjectInfo struct {
	Name           string   `json:"name"`
	Schemes        []string `json:"schemes"`
	Configurations []string `json:"configurations"`
}

type WorkspaceInfo struct {
	Name    string   `json:"name"`
	Schemes []string `json:"schemes"`
}

//-----------------------------------------------------------------------------

func DetectXcodeInfo(basePath, fileName string) (*XcodeInfo, error) {
	xi := &XcodeInfo{}

	if err := xi.populate(basePath, fileName); err != nil {
		return nil, err
	}

	return xi, nil
}

//-----------------------------------------------------------------------------

func (xi *XcodeInfo) Configurations() []string {
	if xi.Project != nil {
		return xi.Project.Configurations
	}

	return nil
}

func (xi *XcodeInfo) Name() string {
	if xi.Project != nil {
		return xi.Project.Name
	}

	if xi.Workspace != nil {
		return xi.Workspace.Name
	}

	return ""
}

func (xi *XcodeInfo) Schemes() []string {
	if xi.Project != nil {
		return xi.Project.Schemes
	}

	if xi.Workspace != nil {
		return xi.Workspace.Schemes
	}

	return nil
}

//-----------------------------------------------------------------------------

func (xi *XcodeInfo) populate(basePath, fileName string) error {
	if strings.HasSuffix(fileName, ".xcworkspace") {
		return xi.populateFromWorkspace(basePath, fileName)
	}

	return xi.populateFromProject(basePath, fileName)
}

func (xi *XcodeInfo) populateFromProject(basePath, project string) error {
	task := lib.NewTask("xcodebuild", "-list", "-json", "-project", project)

	task.Cwd = basePath

	jsonString, _, err := task.Run()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonString), xi)
}

func (xi *XcodeInfo) populateFromWorkspace(basePath, workspace string) error {
	task := lib.NewTask("xcodebuild", "-list", "-json", "-workspace", workspace)

	task.Cwd = basePath

	jsonString, _, err := task.Run()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonString), xi)
}
