package util

import (
	"sort"
)

type sortType struct {
	length int
	less   func(int, int) bool
	swap   func(int, int)
}

func (s *sortType) Len() int {
	return s.length
}

func (s *sortType) Less(i, j int) bool {
	return s.less(i, j)
}

func (s *sortType) Swap(i, j int) {
	s.swap(i, j)
}

// Sort orders elements generically. Len specifies the number of items to be
// sorted. Less is a callback that returns true if the item indexed by the
// first integer is less than the item indexed by the second integer. Swap
// exchanges the items identified by the two index values.
func Sort(Len int, Less func(int, int) bool, Swap func(int, int)) {
	sort.Sort(&sortType{length: Len, less: Less, swap: Swap})
}
