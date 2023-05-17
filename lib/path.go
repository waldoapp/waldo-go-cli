package lib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, fs.ErrNotExist)
}

func FindDirectoryPathsMatching(pattern string) []string {
	var result []string

	matches, err := filepath.Glob(pattern)

	if err != nil || len(matches) == 0 {
		return result
	}

	for _, match := range matches {
		if IsDirectory(match) {
			result = append(result, match)
		}
	}

	return result
}

func HasDirectoryMatching(pattern string) bool {
	matches, _ := filepath.Glob(pattern)

	for _, match := range matches {
		if IsDirectory(match) {
			return true
		}
	}

	return false
}

func HasRegularFileMatching(pattern string) bool {
	matches, _ := filepath.Glob(pattern)

	for _, match := range matches {
		if IsRegularFile(match) {
			return true
		}
	}

	return false
}

func IsDirectory(path string) bool {
	fi, err := os.Stat(path)

	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return fi.Mode().IsDir()
}

func IsRegularFile(path string) bool {
	fi, err := os.Stat(path)

	if errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return fi.Mode().IsRegular()
}

func MakeRelativeToCWD(path string) string {
	cwd, err := os.Getwd()

	if err == nil {
		relPath, err := filepath.Rel(cwd, path)

		if err == nil {
			return relPath
		}

	}

	return path
}
