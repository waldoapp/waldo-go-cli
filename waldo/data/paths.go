package data

import (
	"os"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func FindRepoSpecificPath() (string, error) {
	path, err := lib.FindGitRepositoryPath()

	if err == nil {
		path = filepath.Join(filepath.Dir(path), ".waldo")
	}

	return path, err
}

func FindUserSpecificPath() (string, error) {
	path, err := os.UserConfigDir()

	if err == nil {
		path, err = filepath.Abs(filepath.Join(path, "waldo"))
	}

	if err != nil {
		return "", err
	}

	return path, nil
}
