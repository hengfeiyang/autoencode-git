package util

import (
	"sort"
)

type sortBy func(p1, p2 map[string]string) bool

type mapSorter struct {
	items []map[string]string
	by    sortBy
}

func (s *mapSorter) Len() int {
	return len(s.items)
}

func (s *mapSorter) Less(p1, p2 int) bool {
	return s.by(s.items[p1], s.items[p2])
}

func (s *mapSorter) Swap(p1, p2 int) {
	s.items[p1], s.items[p2] = s.items[p2], s.items[p1]
}

func MapSort(mapItems []map[string]string, mapBy sortBy) {
	psort := &mapSorter{
		items: mapItems,
		by:    mapBy,
	}
	sort.Sort(psort)
}
