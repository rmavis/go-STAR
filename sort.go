package main





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
