package filter

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type FileExtFilter struct {
	Exts []string
}

func (filter *FileExtFilter) supports(info fs.FileInfo, path string) bool {
	if len(filter.Exts) == 0 {
		return true
	}
	fileExt := filepath.Ext(info.Name())
	if info.Name() != "" {
		fileExt = fileExt[1:]
	}
	for _, ext := range filter.Exts {
		if ext == strings.ToLower(fileExt) {
			return true
		}
	}
	return false
}
