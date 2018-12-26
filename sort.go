package main

import (
	"sort"
)


type ByMatchRate []Record

func (a ByMatchRate) Len() int {
	return len(a)
}

func (a ByMatchRate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByMatchRate) Less(i, j int) bool {
	return a[i].MatchRate < a[j].MatchRate
}


// The ByMatchRate and ByDateCreated types and methods are used
// to sort Records by the specified criteria.


type ByDateCreated []Record

func (a ByDateCreated) Len() int {
	return len(a)
}

func (a ByDateCreated) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDateCreated) Less(i, j int) bool {
	return a[i].Meta[0] < a[j].Meta[0]
}


// makeSorter returns the sorting function used in the multi-part
// Search action function. If search terms are given, then the sort
// sort will be by relevancy. Else, by date.
func makeSorter(act *ActionCode, has_terms bool) func([]Record) {
	var sorter func([]Record)

	if (has_terms) {
		if act.Sort == SortAsc {  // ascending
			sorter = func(records []Record) {
				sort.Sort(ByMatchRate(records))
			}
		} else {  // descending
			sorter = func(records []Record) {
				sort.Sort(sort.Reverse(ByMatchRate(records)))
			}
		}
	} else {
		if act.Sort == SortAsc {  // ascending
			sorter = func(records []Record) {
				sort.Sort(ByDateCreated(records))
			}
		} else {  // descending
			sorter = func(records []Record) {
				sort.Sort(sort.Reverse(ByDateCreated(records)))
			}
		}
	}

	return sorter
}
