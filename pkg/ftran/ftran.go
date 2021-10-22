package ftran

import (
	"github.com/ujum/ftran/internal/filter"
	"github.com/ujum/ftran/internal/transfer"
)

type Options struct {
	SameExtDir   bool
	SourceDir    string
	TargetDir    string
	AffectedExts []*ResourceFilterOption
	AffectedDirs []*ResourceFilterOption
}

type ResourceFilterOption struct {
	Inverse   bool
	Resources []string
}

func Run(opt *Options, resultLogs chan *transfer.ResourceLog) error {
	config := &transfer.Config{
		SameExtDir:    opt.SameExtDir,
		SourceDir:     opt.SourceDir,
		TargetDir:     opt.TargetDir,
		FileFilterReg: createFileFilterRegistry(opt.AffectedExts),
		DirFilterReg:  createDirFilterRegistry(opt.AffectedDirs),
	}
	if err := transfer.Transfer(config, resultLogs); err != nil {
		return err
	}
	return nil
}

func createDirFilterRegistry(dirPathFilterOpt []*ResourceFilterOption) *filter.Registry {
	var dirFilterRegistry = filter.NewFilterRegistry()
	if dirPathFilterOpt != nil {
		addDirNameFilters(dirPathFilterOpt, dirFilterRegistry)
	}
	return dirFilterRegistry
}

func createFileFilterRegistry(fileExtFilterOpt []*ResourceFilterOption) *filter.Registry {
	var fileFilterRegistry = filter.NewFilterRegistry()
	if fileExtFilterOpt != nil {
		addExtFileFilters(fileExtFilterOpt, fileFilterRegistry)
	}
	return fileFilterRegistry
}

func addDirNameFilters(dirPathFilterOpt []*ResourceFilterOption, dirFilterRegistry *filter.Registry) {
	for _, ft := range dirPathFilterOpt {
		var dirNameFilter filter.ResourceFilter = &filter.DirPathFilter{
			Paths: ft.Resources,
		}
		dirNameFilter = makeInverseIfNeed(ft, dirNameFilter)
		dirFilterRegistry.AddFilter(dirNameFilter)
	}
}

func addExtFileFilters(fileExtFilterOpt []*ResourceFilterOption, fileFilterRegistry *filter.Registry) {
	for _, ft := range fileExtFilterOpt {
		var extFileFilter filter.ResourceFilter = &filter.FileExtFilter{
			Exts: ft.Resources,
		}
		extFileFilter = makeInverseIfNeed(ft, extFileFilter)
		fileFilterRegistry.AddFilter(extFileFilter)
	}
}

func makeInverseIfNeed(filterOpt *ResourceFilterOption, resourceFilter filter.ResourceFilter) filter.ResourceFilter {
	if filterOpt.Inverse {
		resourceFilter = filter.Inverse(resourceFilter)
	}
	return resourceFilter
}
