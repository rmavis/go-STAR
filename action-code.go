package main

type ActionCode struct {
	Main int
	Sub int
	Match int
	Sort int
	Print int
}

const (
	MainActView int = iota + 1
	MainActCreate
	MainActHelp
	MainActInit
	MainActDemo
)

const (
	SubActConfig int = iota
	SubActView
	SubActPipe
	SubActEdit
	SubActDelete
)

const (
	MatchConfig int = iota
	MatchLoose
	MatchStrict
)

const (
	SortConfig int = iota
	SortDesc
	SortAsc
)

const (
	PrintConfig int = iota
	PrintFull
	PrintCompact
)
