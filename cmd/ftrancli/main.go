package main

import (
	"flag"
	"fmt"
	"github.com/ujum/ftran/internal/app"
	"os"
	"strings"
)

const (
	paramSeparator = ","
	reversePrefix  = "!"
)

func main() {
	sameExtDir := flag.Bool("oneDir", true, "Flag to move files with the same extensions to one dir")
	targetDir := flag.String("targetDir", "result", "Target directory name")
	affectedExts := flag.String("exts", "", "Restrict a number of affected file extensions (empty string - will affect all extensions).\n"+
		"The name of an extension must be separated by a comma.\n"+
		"Use '"+reversePrefix+"' prefix for reverse")
	affectedDirs := flag.String("dirs", "", "Restrict a number of affected directories (relative path).\n"+
		"Use '"+reversePrefix+"' prefix for reverse")
	flag.Parse()

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("can't get work directory: %v", err)
		return
	}

	dirPathFilterOpt := createResourceFilterOpt(*affectedDirs)
	fileExtFilterOpt := createResourceFilterOpt(strings.ReplaceAll(*affectedExts, ".", ""))
	err = app.Run(&app.Options{
		SameExtDir:   *sameExtDir,
		WorkDir:      workDir,
		TargetDir:    *targetDir,
		AffectedExts: fileExtFilterOpt,
		AffectedDirs: dirPathFilterOpt,
	})
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func createResourceFilterOpt(affectedRes string) *app.ResourceFilterOption {
	if affectedRes != "" {
		inverse := false
		if strings.HasPrefix(affectedRes, reversePrefix) {
			inverse = true
			affectedRes = affectedRes[len(reversePrefix):]
		}
		return &app.ResourceFilterOption{
			Inverse:   inverse,
			Resources: strings.Split(strings.ToLower(affectedRes), paramSeparator),
		}
	}
	return nil
}
