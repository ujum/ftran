package main

import (
	"flag"
	"github.com/ujum/ftran/pkg"
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
	dirPathFilterOpt := createResourceFilterOpt(*affectedDirs)
	fileExtFilterOpt := createResourceFilterOpt(strings.ReplaceAll(*affectedExts, ".", ""))
	err = pkg.Run(&pkg.Options{
		SameExtDir:   *sameExtDir,
		SourceDir:    workDir,
		TargetDir:    *targetDir,
		AffectedExts: fileExtFilterOpt,
		AffectedDirs: dirPathFilterOpt,
	})
	if err != nil {
		log.Printf("error: %v", err)
	}
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

func createResourceFilterOpt(affectedRes string) *pkg.ResourceFilterOption {
	if affectedRes != "" {
		inverse := false
		if strings.HasPrefix(affectedRes, reversePrefix) {
			inverse = true
			affectedRes = affectedRes[len(reversePrefix):]
		}
		return &pkg.ResourceFilterOption{
			Inverse:   inverse,
			Resources: strings.Split(strings.ToLower(affectedRes), paramSeparator),
		}
	}
	return nil
}
