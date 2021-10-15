package transfer

import (
	"fmt"
	"github.com/ujum/ftran/internal/filter"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	SameExtDir    bool
	SourceDir     string
	TargetDir     string
	FileFilterReg *filter.Registry
	DirFilterReg  *filter.Registry
}

func Transfer(config *Config) error {
	err := createDirIfNotExist(config.TargetDir)
	if err != nil {
		return err
	}
	return walkAndMove(config.SourceDir, config.TargetDir, config.SameExtDir, config.FileFilterReg, config.DirFilterReg)
}

func walkAndMove(sourceDir, targetDir string, sameExtDir bool, fileFilterReg *filter.Registry,
	dirFilterReg *filter.Registry) error {
	fileInfo, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	for _, file := range fileInfo {
		processResource(sourceDir, targetDir, sameExtDir, file, fileFilterReg, dirFilterReg)
	}
	return nil
}

func processResource(sourceDir, targetDir string, sameExtDir bool, file fs.FileInfo,
	fileFilterReg *filter.Registry, dirFilterReg *filter.Registry) {
	if file.IsDir() {
		processDir(sourceDir, targetDir, sameExtDir, file, fileFilterReg, dirFilterReg)
	} else {
		processFile(sourceDir, targetDir, fileFilterReg, file)
	}
}

func processFile(sourceDir string, targetDir string, fileFilterReg *filter.Registry, file fs.FileInfo) {
	fileName := file.Name()
	if fileFilterReg.Apply(file, sourceDir) {
		if err := moveFileToExtDir(sourceDir, targetDir, fileName); err != nil {
			fmt.Printf("can't move file [%s]: %v", fileName, err)
		}
	} else {
		fmt.Printf("skipped file: %s", filepath.Join(sourceDir, fileName))
	}
}

func processDir(sourceDir string, targetDir string, sameExtDir bool, file fs.FileInfo,
	fileFilterReg *filter.Registry, dirFilterReg *filter.Registry) {
	nestedTargetDir, err := getOrCreateNestedTargetDir(targetDir, sameExtDir, file.Name())
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	nestedSourceDir := filepath.Join(sourceDir, file.Name())
	if dirFilterReg.Apply(file, sourceDir) {
		if err := walkAndMove(nestedSourceDir, nestedTargetDir, sameExtDir, fileFilterReg, dirFilterReg); err != nil {
			fmt.Printf("can't read directory %s: %v", sourceDir, err)
		}
	} else {
		fmt.Printf("skipped directory: %s", nestedSourceDir)
	}
}

func getOrCreateNestedTargetDir(targetDir string, sameExtDir bool, fileName string) (string, error) {
	toDirNested := targetDir
	if !sameExtDir {
		toDirNested = filepath.Join(targetDir, fileName)
		err := createDirIfNotExist(toDirNested)
		if err != nil {
			return toDirNested, fmt.Errorf("can't create directory [%s]: %v", toDirNested, err)
		}
	}
	return toDirNested, nil
}
