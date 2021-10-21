package pkg

import (
	"github.com/ujum/ftran/internal/filter"
	"github.com/ujum/ftran/internal/transfer"
	"path/filepath"
)

type Options struct {
	SameExtDir   bool
	SourceDir    string
	TargetDir    string
	AffectedExts *ResourceFilterOption
	AffectedDirs *ResourceFilterOption
}

type ResourceFilterOption struct {
	Inverse   bool
	Resources []string
}

func Run(opt *Options) error {
	config := &transfer.Config{
		SameExtDir:    opt.SameExtDir,
		SourceDir:     opt.SourceDir,
		TargetDir:     filepath.Join(filepath.Dir(opt.SourceDir), opt.TargetDir),
		FileFilterReg: createFileFilterRegistry(opt.AffectedExts),
		DirFilterReg:  createDirFilterRegistry(opt.AffectedDirs),
	}
	if err := transfer.Transfer(config); err != nil {
		return err
	}
	return nil
}

func createDirFilterRegistry(dirPathFilterOpt *ResourceFilterOption) *filter.Registry {
	var dirFilterRegistry = filter.NewFilterRegistry()
	if dirPathFilterOpt != nil {
		addDirNameFilter(dirPathFilterOpt, dirFilterRegistry)
	}
	return dirFilterRegistry
}

func createFileFilterRegistry(fileExtFilterOpt *ResourceFilterOption) *filter.Registry {
	var fileFilterRegistry = filter.NewFilterRegistry()
	if fileExtFilterOpt != nil {
		addExtFileFilter(fileExtFilterOpt, fileFilterRegistry)
	}
	return fileFilterRegistry
}

func addDirNameFilter(dirPathFilterOpt *ResourceFilterOption, dirFilterRegistry *filter.Registry) {
	var dirNameFilter filter.ResourceFilter = &filter.DirPathFilter{
		Paths: dirPathFilterOpt.Resources,
	}
	dirNameFilter = makeInverseIfNeed(dirPathFilterOpt, dirNameFilter)
	dirFilterRegistry.AddFilter(dirNameFilter)
}

func addExtFileFilter(fileExtFilterOpt *ResourceFilterOption, fileFilterRegistry *filter.Registry) {
	var extFileFilter filter.ResourceFilter = &filter.FileExtFilter{
		Exts: fileExtFilterOpt.Resources,
	}
	extFileFilter = makeInverseIfNeed(fileExtFilterOpt, extFileFilter)
	fileFilterRegistry.AddFilter(extFileFilter)
}

func makeInverseIfNeed(filterOpt *ResourceFilterOption, resourceFilter filter.ResourceFilter) filter.ResourceFilter {
	if filterOpt.Inverse {
		resourceFilter = filter.Inverse(resourceFilter)
	}
	return resourceFilter
}
