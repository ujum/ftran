package transfer

import (
	"fmt"
	"github.com/ujum/ftran/internal/filter"
	"github.com/ujum/ftran/pkg/data"
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

func Transfer(config *Config, resourceLogs chan *data.ResourceLog) error {
	defer close(resourceLogs)
	err := createDirIfNotExist(config.TargetDir)
	if err != nil {
		return err
	}
	return walkAndMove(config, resourceLogs)
}

func walkAndMove(config *Config, resourceLogs chan *data.ResourceLog) error {
	fileInfo, err := ioutil.ReadDir(config.SourceDir)
	if err != nil {
		return err
	}
	for _, file := range fileInfo {
		processResource(config, file, resourceLogs)
	}
	return nil
}

func processResource(config *Config, file fs.FileInfo, resourceLogs chan *data.ResourceLog) {
	if file.IsDir() {
		processDir(*config, file, resourceLogs)
	} else {
		processFile(config, file, resourceLogs)
	}
}

func processFile(config *Config, file fs.FileInfo, resourceLogs chan *data.ResourceLog) {
	fileName := file.Name()
	source := filepath.Join(config.SourceDir, fileName)
	target := filepath.Join(config.TargetDir, fileName)
	log := data.NewResourceLog(source, target, false, nil)
	if config.FileFilterReg == nil || config.FileFilterReg.Apply(file, config.SourceDir) {
		if err := moveFileToExtDir(config.SourceDir, config.TargetDir, fileName); err != nil {
			log.Error = err
			log.Skipped = true
		}
	} else {
		log.Skipped = true
	}
	resourceLogs <- log
}

func processDir(config Config, file fs.FileInfo, resourceLogs chan *data.ResourceLog) {
	nestedSourceDir := filepath.Join(config.SourceDir, file.Name())
	nestedTargetDir, err := getOrCreateNestedTargetDir(config.TargetDir, config.SameExtDir, file.Name())
	if err != nil {
		resourceLogs <- data.NewResourceLog(nestedSourceDir, nestedTargetDir, true, err)
		return
	}
	if config.DirFilterReg == nil || config.DirFilterReg.Apply(file, config.SourceDir) {
		config.TargetDir = nestedTargetDir
		config.SourceDir = nestedSourceDir
		if err := walkAndMove(&config, resourceLogs); err != nil {
			resourceLogs <- data.NewResourceLog(nestedSourceDir, nestedTargetDir, true, err)
		}
	} else {
		resourceLogs <- data.NewResourceLog(nestedSourceDir, nestedTargetDir, true, nil)
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
