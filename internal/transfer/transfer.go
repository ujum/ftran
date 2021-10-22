package transfer

import (
	"fmt"
	"github.com/ujum/ftran/internal/filter"
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

type ResourceLog struct {
	Source  string
	Target  string
	Skipped bool
	Error   error
}

type Config struct {
	SameExtDir    bool
	SourceDir     string
	TargetDir     string
	FileFilterReg *filter.Registry
	DirFilterReg  *filter.Registry
}

func newResourceLog(nestedSourceDir string, nestedTargetDir string, skipped bool, err error) *ResourceLog {
	return &ResourceLog{
		Source:  nestedSourceDir,
		Target:  nestedTargetDir,
		Skipped: skipped,
		Error:   err,
	}
}

func Transfer(config *Config, resourceLogs chan *ResourceLog) error {
	defer close(resourceLogs)
	err := createDirIfNotExist(config.TargetDir)
	if err != nil {
		return err
	}
	return walkAndMove(config, resourceLogs)
}

func walkAndMove(config *Config, resourceLogs chan *ResourceLog) error {
	fileInfo, err := ioutil.ReadDir(config.SourceDir)
	if err != nil {
		return err
	}
	for _, file := range fileInfo {
		processResource(config, file, resourceLogs)
	}
	return nil
}

func processResource(config *Config, file fs.FileInfo, resourceLogs chan *ResourceLog) {
	if file.IsDir() {
		processDir(*config, file, resourceLogs)
	} else {
		processFile(config, file, resourceLogs)
	}
}

func processFile(config *Config, file fs.FileInfo, resourceLogs chan *ResourceLog) {
	fileName := file.Name()
	source := filepath.Join(config.SourceDir, fileName)
	target := filepath.Join(config.TargetDir, fileName)
	if config.FileFilterReg == nil || config.FileFilterReg.Apply(file, config.SourceDir) {
		if err := moveFileToExtDir(config.SourceDir, config.TargetDir, fileName); err != nil {
			resourceLogs <- newResourceLog(source, target, true, err)
		}
	} else {
		resourceLogs <- newResourceLog(source, target, true, nil)
	}
}

func processDir(config Config, file fs.FileInfo, resourceLogs chan *ResourceLog) {
	nestedSourceDir := filepath.Join(config.SourceDir, file.Name())
	nestedTargetDir, err := getOrCreateNestedTargetDir(config.TargetDir, config.SameExtDir, file.Name())
	if err != nil {
		resourceLogs <- newResourceLog(nestedSourceDir, nestedTargetDir, true, err)
		return
	}
	if config.DirFilterReg == nil || config.DirFilterReg.Apply(file, config.SourceDir) {
		config.TargetDir = nestedTargetDir
		config.SourceDir = nestedSourceDir
		if err := walkAndMove(&config, resourceLogs); err != nil {
			resourceLogs <- newResourceLog(nestedSourceDir, nestedTargetDir, true, err)
		}
	} else {
		resourceLogs <- newResourceLog(nestedSourceDir, nestedTargetDir, true, nil)
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
