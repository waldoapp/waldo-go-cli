package lib

import (
	"errors"
	"os"
	"path/filepath"
)

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
