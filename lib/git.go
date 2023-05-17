package lib

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type GitInfo struct {
	Access GitAccess
	Branch string
	Commit string
}

//-----------------------------------------------------------------------------

type GitAccess int

const (
	GitAccessOk GitAccess = iota + 1 // MUST be first
	GitAccessNoGitCommandFound
	GitAccessNotGitRepository
)

//-----------------------------------------------------------------------------

func DetectGitAccess() GitAccess {
	access := GitAccessOk

	if !isGitInstalled() {
		access = GitAccessNoGitCommandFound
	} else if !hasGitRepository() {
		access = GitAccessNotGitRepository
	}

	return access
}

func FindGitRepositoryPath() (string, error) {
	cfgName := ".git/config"
	repPath := ""

	dirPath, err := os.Getwd()

	if err != nil {
		return "", err
	}

	for {
		tstPath := filepath.Join(dirPath, cfgName)

		if IsRegularFile(tstPath) {
			repPath = filepath.Dir(tstPath)
			break
		}

		tmpPath := filepath.Dir(dirPath)

		if tmpPath == dirPath {
			return "", errors.New("Not a git repository (or any of the parent directories)")
		}

		dirPath = tmpPath
	}

	return filepath.Abs(repPath)
}

func InferGitInfo(skipCount int) *GitInfo {
	access := GitAccessOk
	branch := ""
	commit := ""

	if !isGitInstalled() {
		access = GitAccessNoGitCommandFound
	} else if !hasGitRepository() {
		access = GitAccessNotGitRepository
	} else {
		commit = inferGitCommit(skipCount)
		branch = inferGitBranch(commit)
	}

	return &GitInfo{
		Access: access,
		Branch: branch,
		Commit: commit}
}

//-----------------------------------------------------------------------------

func (ga GitAccess) String() string {
	return [...]string{
		"ok",
		"noGitCommandFound",
		"notGitRepository"}[ga-1]
}

func fetchBranchNamesFromGitForEachRefResults(results string) []string {
	lines := strings.Split(results, "\n")

	var branchNames []string

	for _, line := range lines {
		branchName := refNameToBranchName(line)

		if len(branchName) > 0 {
			branchNames = append(branchNames, branchName)
		}
	}

	return removeDuplicates(branchNames)
}

func hasGitRepository() bool {
	_, _, err := NewTask("git", []string{"rev-parse"}).Run()

	return err == nil
}

func inferGitBranch(commit string) string {
	if len(commit) > 0 {
		fromForEachRev := inferGitBranchFromForEachRef(commit)

		if len(fromForEachRev) > 0 {
			return fromForEachRev
		}

		fromNameRev := inferGitBranchFromNameRev(commit)

		if len(fromNameRev) > 0 {
			return fromNameRev
		}
	}

	return inferGitBranchFromRevParse()
}

func inferGitBranchFromForEachRef(commit string) string {
	args := []string{"for-each-ref", fmt.Sprintf("--points-at=%s", commit), "--format=%(refname)"}

	stdout, _, err := NewTask("git", args).Run()

	if err == nil {
		branchNames := fetchBranchNamesFromGitForEachRefResults(stdout)

		if len(branchNames) > 0 {
			//
			// Since we don’t know which branch is the correct one, arbitrarily
			// return the first one:
			//
			return branchNames[0]
		}
	}

	return ""
}

func inferGitBranchFromNameRev(commit string) string {
	args := []string{"name-rev", "--always", "--name-only", commit}

	name, _, err := NewTask("git", args).Run()

	if err == nil {
		return nameRevToBranchName(name)
	}

	return ""
}

func inferGitBranchFromRevParse() string {
	args := []string{"rev-parse", "--abbrev-ref", "HEAD"}

	name, _, err := NewTask("git", args).Run()

	if err == nil && name != "HEAD" {
		return name
	}

	return ""
}

func inferGitCommit(skipCount int) string {
	args := []string{"log", "--format=%H", fmt.Sprintf("--skip=%d", skipCount), "-1"}

	hash, _, err := NewTask("git", args).Run()

	if err != nil {
		return ""
	}

	return hash
}

func isGitInstalled() bool {
	var name string

	if runtime.GOOS == "windows" {
		name = "git.exe"
	} else {
		name = "git"
	}

	_, err := exec.LookPath(name)

	return err == nil
}

func nameRevToBranchName(refName string) string {
	branchName := strings.TrimSpace(refName)

	if strings.HasPrefix(branchName, "tags/") {
		return ""
	}

	if strings.HasPrefix(branchName, "remotes/") {
		branchName = strings.TrimPrefix(branchName, "remotes/")

		//
		// Remove the remote name:
		//
		slash := strings.Index(branchName, "/")

		if slash == -1 {
			return ""
		}

		branchName = branchName[slash+1:]
	}

	if branchName == "HEAD" {
		return ""
	}

	return branchName
}

func refNameToBranchName(refName string) string {
	branchName := strings.TrimSpace(refName)

	if strings.HasPrefix(branchName, "refs/heads/") {
		branchName = strings.TrimPrefix(branchName, "refs/heads/")
	} else if strings.HasPrefix(branchName, "refs/remotes/") {
		branchName = strings.TrimPrefix(branchName, "refs/remotes/")

		//
		// Remove the remote name:
		//
		slash := strings.Index(branchName, "/")

		if slash == -1 {
			return ""
		}

		branchName = branchName[slash+1:]
	} else {
		return ""
	}

	if branchName == "HEAD" {
		return ""
	}

	return branchName
}

func removeDuplicates(strings []string) []string {
	seen := make(map[string]bool)

	var result []string

	for _, s := range strings {
		if _, ok := seen[s]; !ok {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}
