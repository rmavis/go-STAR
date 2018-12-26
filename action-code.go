package main


// ActionCode is a structure that contains fields for encoding the
// user's desired action as parsed from the command line.
type ActionCode struct {
	Main int
	Sub int
	Match int
	Sort int
	Print int
}

// These constants are like enums. They clarify the purpose of an
// ActionCode's fields.

const (
	MainActView int = iota
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


// mergeConfigActions receives pointers to a Config and an ActionCode
// and sets values in the ActionCode according to values in the Config.
func mergeConfigActions(conf *Config, act *ActionCode) {
	if act.Sub == SubActConfig {
		if len(conf.Action) > 0 {
			act.Sub = SubActPipe
		} else {
			act.Sub = SubActView
		}
	}

	if act.Match == MatchConfig {
		if conf.FilterMode == "strict" {
			act.Match = MatchStrict
		} else {
			act.Match = MatchLoose
		}
	}

	if act.Sort == SortConfig {
		if conf.SortOrder == "asc" {
			act.Sort = SortAsc
		} else {
			act.Sort = SortDesc
		}
	}

	if act.Print == PrintConfig {
		if conf.PrintLines == "2" {
			act.Print = PrintFull
		} else {
			act.Print = PrintCompact
		}
	}
}
