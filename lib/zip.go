package lib

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func ZipFolder(zipPath, folderPath, basePath string) error {
	err := os.Chdir(folderPath)

	if err != nil {
		return err
	}

	zipFile, err := os.Create(zipPath)

	if err != nil {
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	walker := func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		file, err := os.Open(path)

		if err != nil {
			return err
		}

		defer file.Close()

		zipEntry, err := zipWriter.Create(path)

		if err != nil {
			return err
		}

		_, err = io.Copy(zipEntry, file)

		return err
	}

	err = filepath.WalkDir(basePath, walker)

	err2 := zipWriter.Close()

	if err != nil {
		return err
	}

	return err2
}
