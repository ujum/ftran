package transfer

import (
	"fmt"
	"github.com/ujum/ftran/internal/filter"
	"io/fs"
	"io/ioutil"
	"log"
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
	return walkAndMove(config)
}

func walkAndMove(config *Config) error {
	fileInfo, err := ioutil.ReadDir(config.SourceDir)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	for _, file := range fileInfo {
		processResource(config, file)
	}
	return nil
}

func processResource(config *Config, file fs.FileInfo) {
	if file.IsDir() {
		processDir(*config, file)
	} else {
		processFile(config, file)
	}
}

func processFile(config *Config, file fs.FileInfo) {
	fileName := file.Name()
	if config.FileFilterReg == nil || config.FileFilterReg.Apply(file, config.SourceDir) {
		if err := moveFileToExtDir(config.SourceDir, config.TargetDir, fileName); err != nil {
			log.Printf("can't move file [%s]: %v", fileName, err)
		}
	} else {
		log.Printf("skipped file: %s", filepath.Join(config.SourceDir, fileName))
	}
}

func processDir(config Config, file fs.FileInfo) {
	nestedTargetDir, err := getOrCreateNestedTargetDir(config.TargetDir, config.SameExtDir, file.Name())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	nestedSourceDir := filepath.Join(config.SourceDir, file.Name())
	if config.DirFilterReg == nil || config.DirFilterReg.Apply(file, config.SourceDir) {
		config.TargetDir = nestedTargetDir
		config.SourceDir = nestedSourceDir
		if err := walkAndMove(&config); err != nil {
			log.Printf("can't read directory %s: %v", config.SourceDir, err)
		}
	} else {
		log.Printf("skipped directory: %s", nestedSourceDir)
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
