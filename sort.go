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





func makeSorter(conf Config, action_code []int, terms []string) func([]Record) {
	var sorter func([]Record)

	switch {
	case action_code[0] == 4:
		sorter = func(records []Record) {
			sort.Sort(ByDateCreated(records))
		}
	case len(terms) == 0 || conf.Action == "browse" || action_code[1] == 5:
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
