package main

import (
	"flag"
	"github.com/ujum/ftran/pkg/data"
	"github.com/ujum/ftran/pkg/ftran"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	paramSeparator = ","
	reversePrefix  = "!"
)

func main() {
	sameExtDir := flag.Bool("oneDir", true, "Flag to move files with the same extensions to one dir")
	sourceDir := flag.String("sourceDir", "", "Source directory.")
	targetDir := flag.String("targetDir", "result", "Target directory name")
	affectedExts := flag.String("exts", "", "Restrict a number of affected file extensions (empty string - will affect all extensions).\n"+
		"The name of an extension must be separated by a comma.\n"+
		"Use '"+reversePrefix+"' prefix for reverse")
	affectedDirs := flag.String("dirs", "", "Restrict a number of affected directories (relative path).\n"+
		"Use '"+reversePrefix+"' prefix for reverse")
	flag.Parse()

	if *sourceDir == "" {
		log.Println("-sourceDir argument not provided")
		flag.Usage()
		return
	}

	workDir, err := getWorkDir(*sourceDir)
	if err != nil {
		log.Printf("can't get work directory: %v", err)
		return
	}
	var allAffectedExts []*ftran.ResourceFilterOption
	var allAffectedDirs []*ftran.ResourceFilterOption
	dirPathFilterOpt := createResourceFilterOpt(*affectedDirs)
	if dirPathFilterOpt != nil {
		allAffectedDirs = append(allAffectedDirs, dirPathFilterOpt)
	}
	fileExtFilterOpt := createResourceFilterOpt(strings.ReplaceAll(*affectedExts, ".", ""))
	if fileExtFilterOpt != nil {
		allAffectedExts = append(allAffectedExts, fileExtFilterOpt)
	}
	resourceLogs := make(chan *data.ResourceLog)
	go func() {
		err = ftran.Run(&ftran.Options{
			SameExtDir:   *sameExtDir,
			SourceDir:    workDir,
			TargetDir:    filepath.Join(filepath.Dir(workDir), *targetDir) + "_" + filepath.Base(workDir),
			AffectedExts: allAffectedExts,
			AffectedDirs: allAffectedDirs,
		}, resourceLogs)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()
	printResourceLogs(resourceLogs)
}

func getWorkDir(sourceDir string) (string, error) {
	if filepath.IsAbs(sourceDir) {
		return sourceDir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, sourceDir), nil
}

func createResourceFilterOpt(affectedRes string) *ftran.ResourceFilterOption {
	if affectedRes != "" {
		inverse := false
		if strings.HasPrefix(affectedRes, reversePrefix) {
			inverse = true
			affectedRes = affectedRes[len(reversePrefix):]
		}
		return &ftran.ResourceFilterOption{
			Inverse:   inverse,
			Resources: strings.Split(strings.ToLower(affectedRes), paramSeparator),
		}
	}
	return nil
}

func printResourceLogs(result chan *data.ResourceLog) {
	for res := range result {
		if res.Skipped {
			if res.Error == nil {
				log.Printf("skipped: %s\n", res.Source)
			} else {
				log.Printf("skipped %s\n resource cause: [%v]", res.Source, res.Error)
			}
		} else {
			log.Printf("%s --> %s\n", res.Source, res.Target)
		}
	}
}
