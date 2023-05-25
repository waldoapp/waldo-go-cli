package data

import (
	"os"
	"path/filepath"

	"github.com/waldoapp/waldo-go-cli/lib"
)

func FindRepoSpecificPath() (string, error) {
	path, err := lib.FindGitRepositoryPath()

	if err != nil {
		return "", err
	}

	path = filepath.Join(filepath.Dir(path), ".waldo")

	return path, nil
}

func FindUserSpecificPath() (string, error) {
	path, err := os.UserConfigDir()

	if err != nil {
		return "", err
	}

	path = filepath.Join(path, "waldo")

	return path, nil
}
