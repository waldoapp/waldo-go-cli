package lib

import (
	"os"
	"strings"
)

type CIInfo struct {
	GitBranch string
	GitCommit string
	Provider  CIProvider
	SkipCount int
}

//-----------------------------------------------------------------------------

type CIProvider int

const (
	CIProviderUnknown CIProvider = iota // MUST be first
	CIProviderAppCenter
	CIProviderAzureDevOps
	CIProviderBitrise
	CIProviderCircleCI
	CIProviderCodeBuild
	CIProviderGitHubActions
	CIProviderJenkins
	CIProviderTeamCity
	CIProviderTravisCI
	CIProviderXcodeCloud
)

//-----------------------------------------------------------------------------

func DetectCIInfo(fullInfo bool) *CIInfo {
	info := &CIInfo{
		Provider: detectCIProvider()}

	if fullInfo {
		info.extractFullInfo()
	}

	return info
}

//-----------------------------------------------------------------------------

func (cp CIProvider) String() string {
	return [...]string{
		"Unknown",
		"App Center",
		"Azure DevOps",
		"Bitrise",
		"CircleCI",
		"CodeBuild",
		"GitHub Actions",
		"Jenkins",
		"TeamCity",
		"Travis CI",
		"Xcode Cloud"}[cp]
}

//-----------------------------------------------------------------------------

func (ci *CIInfo) extractFullInfo() {
	switch ci.Provider {
	case CIProviderAppCenter:
		ci.extractFullInfoFromAppCenter()

	case CIProviderAzureDevOps:
		ci.extractFullInfoFromAzureDevOps()

	case CIProviderBitrise:
		ci.extractFullInfoFromBitrise()

	case CIProviderCircleCI:
		ci.extractFullInfoFromCircleCI()

	case CIProviderCodeBuild:
		ci.extractFullInfoFromCodeBuild()

	case CIProviderGitHubActions:
		ci.extractFullInfoFromGitHubActions()

	case CIProviderJenkins:
		ci.extractFullInfoFromJenkins()

	case CIProviderTeamCity:
		ci.extractFullInfoFromTeamCity()

	case CIProviderTravisCI:
		ci.extractFullInfoFromTravisCI()

	case CIProviderXcodeCloud:
		ci.extractFullInfoFromXcodeCloud()

	default:
		break
	}
}

func (ci *CIInfo) extractFullInfoFromAppCenter() {
	//
	// https://docs.microsoft.com/en-us/appcenter/build/custom/variables/
	//
	ci.GitBranch = os.Getenv("APPCENTER_BRANCH")
	ci.GitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromAzureDevOps() {
	//
	// https://docs.microsoft.com/en-us/azure/devops/pipelines/build/variables#build-variables-devops-services
	//
	ci.GitBranch = os.Getenv("BUILD_SOURCEBRANCHNAME")
	ci.GitCommit = os.Getenv("BUILD_SOURCEVERSION")
}

func (ci *CIInfo) extractFullInfoFromBitrise() {
	//
	// https://devcenter.bitrise.io/en/references/available-environment-variables.html
	//
	ci.GitBranch = os.Getenv("BITRISE_GIT_BRANCH")
	ci.GitCommit = os.Getenv("BITRISE_GIT_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromCircleCI() {
	//
	// https://circleci.com/docs2/2.0/env-vars#built-in-environment-variables
	//
	ci.GitBranch = os.Getenv("CIRCLE_BRANCH")
	ci.GitCommit = os.Getenv("CIRCLE_SHA1")
}

func (ci *CIInfo) extractFullInfoFromCodeBuild() {
	//
	// https://docs.aws.amazon.com/codebuild/latest/userguide/build-env-ref-env-vars.html
	//
	trigger := os.Getenv("CODEBUILD_WEBHOOK_TRIGGER")

	if strings.HasPrefix(trigger, "branch/") {
		ci.GitBranch = strings.TrimPrefix(trigger, "branch/")
	} else {
		ci.GitBranch = ""
	}

	ci.GitCommit = os.Getenv("CODEBUILD_WEBHOOK_PREV_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromGitHubActions() {
	//
	// https://docs.github.com/en/actions/learn-github-actions/environment-variables#default-environment-variables
	//
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	refType := os.Getenv("GITHUB_REF_TYPE")

	switch eventName {
	case "pull_request", "pull_request_target":
		if refType == "branch" {
			ci.GitBranch = os.Getenv("GITHUB_HEAD_REF")
		} else {
			ci.GitBranch = ""
		}

		//
		// The following environment variable must be set by us (most likely in
		// a custom action) to match the current value of
		// `github.event.pull_request.head.sha`:
		//
		ci.GitBranch = os.Getenv("GITHUB_EVENT_PULL_REQUEST_HEAD_SHA")

		ci.SkipCount = 1

	case "push":
		if refType == "branch" {
			ci.GitBranch = os.Getenv("GITHUB_REF_NAME")
		} else {
			ci.GitBranch = ""
		}

		ci.GitCommit = os.Getenv("GITHUB_SHA")

	default:
		ci.GitBranch = ""
		ci.GitCommit = ""
	}
}

func (ci *CIInfo) extractFullInfoFromJenkins() {
	ci.GitBranch = "" //os.Getenv("???") -- not currently supported?
	ci.GitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromTeamCity() {
	ci.GitBranch = "" //os.Getenv("???") -- not currently supported?
	ci.GitCommit = "" //os.Getenv("???") -- not currently supported?
}

func (ci *CIInfo) extractFullInfoFromTravisCI() {
	//
	// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
	//
	ci.GitBranch = os.Getenv("TRAVIS_BRANCH")
	ci.GitCommit = os.Getenv("TRAVIS_COMMIT")
}

func (ci *CIInfo) extractFullInfoFromXcodeCloud() {
	//
	// https://developer.apple.com/documentation/xcode/environment-variable-reference
	//
	ci.GitBranch = os.Getenv("CI_BRANCH")

	if ci.GitBranch == "" {
		ci.GitBranch = os.Getenv("CI_PULL_REQUEST_SOURCE_BRANCH")
	}

	ci.GitCommit = os.Getenv("CI_COMMIT")

	if ci.GitCommit == "" {
		ci.GitCommit = os.Getenv("CI_PULL_REQUEST_SOURCE_COMMIT")
	}
}

//-----------------------------------------------------------------------------

func detectCIProvider() CIProvider {
	switch {
	case onAppCenter():
		return CIProviderAppCenter

	case onAzureDevOps():
		return CIProviderAzureDevOps

	case onBitrise():
		return CIProviderBitrise

	case onCircleCI():
		return CIProviderCircleCI

	case onCodeBuild():
		return CIProviderCodeBuild

	case onGitHubActions():
		return CIProviderGitHubActions

	case onJenkins():
		return CIProviderJenkins

	case onTeamCity():
		return CIProviderTeamCity

	case onTravisCI():
		return CIProviderTravisCI

	case onXcodeCloud():
		return CIProviderXcodeCloud

	default:
		return CIProviderUnknown
	}
}

func onAppCenter() bool {
	return len(os.Getenv("APPCENTER_BUILD_ID")) > 0
}

func onAzureDevOps() bool {
	return len(os.Getenv("AGENT_ID")) > 0
}

func onBitrise() bool {
	return os.Getenv("BITRISE_IO") == "true"
}

func onCircleCI() bool {
	return os.Getenv("CIRCLECI") == "true"
}

func onCodeBuild() bool {
	return len(os.Getenv("CODEBUILD_BUILD_ID")) > 0
}

func onGitHubActions() bool {
	return os.Getenv("GITHUB_ACTIONS") == "true"
}

func onJenkins() bool {
	return len(os.Getenv("JENKINS_URL")) > 0
}

func onTeamCity() bool {
	return len(os.Getenv("TEAMCITY_VERSION")) > 0
}

func onTravisCI() bool {
	return os.Getenv("TRAVIS") == "true"
}

func onXcodeCloud() bool {
	return len(os.Getenv("CI_BUILD_ID")) > 0
}
