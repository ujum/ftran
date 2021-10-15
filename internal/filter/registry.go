package filter

import (
	"io/fs"
)

type Registry struct {
	filters []ResourceFilter
}

func NewFilterRegistry(filters ...ResourceFilter) *Registry {
	return &Registry{filters: filters}
}

func (registry *Registry) AddFilter(filter ResourceFilter) {
	registry.filters = append(registry.filters, filter)
}

func (registry *Registry) Apply(info fs.FileInfo, path string) bool {
	if len(registry.filters) == 0 {
		return true
	}
	for _, filter := range registry.filters {
		if !filter.supports(info, path) {
			return false
		}
	}
	return true
}
