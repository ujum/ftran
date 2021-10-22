package transfer

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	perm      = 0755
	extPrefix = "ext_"
	absent    = "absent"
)

func moveFileToExtDir(sourceDir, targetRootDir, file string) error {
	ext, err := createExtDir(targetRootDir, file)
	if err != nil {
		return err
	}
	err = moveFileToDir(sourceDir, ext, file)
	if err != nil {
		return err
	}
	return nil
}

func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, perm)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func createExtDir(targetDir, file string) (string, error) {
	fileExt := filepath.Ext(file)
	if fileExt == "" {
		fileExt = absent
	} else {
		fileExt = fileExt[1:]
	}
	ext := filepath.Join(targetDir, strings.ToUpper(extPrefix+fileExt))
	return ext, createDirIfNotExist(ext)
}

func moveFileToDir(sourceDir, targetDir, file string) error {
	return os.Rename(filepath.Join(sourceDir, file), filepath.Join(targetDir, file))
}
