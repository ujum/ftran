package filter

import (
	"io/fs"
)

type ResourceFilter interface {
	supports(info fs.FileInfo, path string) bool
}

type inverse struct {
	Filter ResourceFilter
}

func Inverse(filter ResourceFilter) ResourceFilter {
	return &inverse{Filter: filter}
}

func (i *inverse) supports(info fs.FileInfo, path string) bool {
	return !i.Filter.supports(info, path)
}
