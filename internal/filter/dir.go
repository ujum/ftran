package filter

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type DirPathFilter struct {
	Paths []string
}

func (filter *DirPathFilter) supports(info fs.FileInfo, path string) bool {
	if len(filter.Paths) == 0 {
		return true
	}
	dirPath := strings.ToLower(filepath.Join(path, info.Name()))
	for _, filterPath := range filter.Paths {
		if strings.Contains(dirPath, filterPath) {
			return true
		}
	}
	return false
}
