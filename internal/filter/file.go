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
	for _, ext := range filter.Exts {
		if ext == strings.ToLower(filepath.Ext(info.Name())[1:]) {
			return true
		}
	}
	return false
}
