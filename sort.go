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
// Search action function.
func makeSorter(action_code []int) func([]Record) {
	var sorter func([]Record)

	switch {
	case action_code[0] == 4:  // Dump in the order created.
		sorter = func(records []Record) {
			sort.Sort(ByDateCreated(records))
		}
	case action_code[1] == 5:  // Browse recent-to-old.
		sorter = func(records []Record) {
			sort.Sort(sort.Reverse(ByDateCreated(records)))
		}
	default:
		sorter = func(records []Record) {
			sort.Sort(sort.Reverse(ByMatchRate(records)))
		}
	}

	return sorter
}
