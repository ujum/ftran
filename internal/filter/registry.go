package filter

import (
	"io/fs"
)

type FilterRegistry struct {
	filters []ResourceFilter
}

func NewFilterRegistry(filters ...ResourceFilter) *FilterRegistry {
	return &FilterRegistry{filters: filters}
}

func (registry *FilterRegistry) AddFilter(filter ResourceFilter) {
	registry.filters = append(registry.filters, filter)
}

func (registry *FilterRegistry) Apply(info fs.FileInfo, path string) bool {
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
