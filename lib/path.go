package lib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
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

func GetModificationTimeUTC(path string) time.Time {
	fi, err := os.Stat(path)

	if err != nil {
		return time.Time{}
	}

	return fi.ModTime().UTC()
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

func MakeRelative(path, basePath string) string {
	relPath, err := filepath.Rel(basePath, path)

	if err != nil {
		return path
	}

	return relPath
}

func MakeRelativeToCWD(path string) string {
	cwd, err := os.Getwd()

	if err != nil {
		return path
	}

	return MakeRelative(path, cwd)
}
