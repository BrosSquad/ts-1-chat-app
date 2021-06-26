package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GetAbsolutePath(path string) (string, error) {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)

		if err != nil {
			return "", err
		}

		return path, nil
	}

	return "", err
}

func CreatePath(path string, perm fs.FileMode) (string, error) {
	path, err := GetAbsolutePath(path)

	if err != nil {
		return "", err
	}

	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, perm); err != nil {
		return "", err
	}

	return path, nil
}

func CreateLogFile(path string, perm fs.FileMode) (file *os.File, err error) {
	path, err = CreatePath(path, perm)

	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(path)

			if err != nil {
				return nil, err
			}

			if err := file.Chmod(perm); err != nil {
				return nil, err
			}

			if err := file.Close(); err != nil {
				return nil, err
			}

		}
	}

	file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND, os.ModeAppend)

	return
}
