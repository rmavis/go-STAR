package main

import (
	"fmt"
	"os"
	"strings"
)


// parseArgs takes a slice of strings, being the command line args,
// and returns an ActionCode and a slice of strings that, if present,
// will affect the action (as search terms, new inputs, etc).
func parseArgs(args []string) (*ActionCode, []string) {
	act := &ActionCode{MainActView, SubActConfig, MatchConfig, SortConfig, PrintConfig}
	var strs []string

out:
	for o := 0; o < len(args); o++ {
		// Command options start with dashes.
		if args[o][0] == '-' {
			arg := strings.ToLower(args[o])

			if arg[1] == '-' {  // Long-form options start with two.
				arg = strings.Replace(arg, string('-'), "", -1)
				updateActionCodeFromWord(arg, act)
			} else {  // Short-form options start with one.
				arg = strings.Replace(arg, string('-'), "", -1)
				for i := 0; i < len(arg); i++ {
					char := string(arg[i])
					updateActionCodeFromChar(char, act)
				}
			}
		} else {
			strs = args[o:]
			// The first non-option argument signals that the rest of
			// the arguments are for input for the action.
			break out;
		}
	}

	return act, strs
}

// updateActionCodeFromChar receives a string, being a short-form
// command-line option, and a pointer to an ActionCode, and it sets
// some value in that ActionCode according to the string.
func updateActionCodeFromChar(arg string, act *ActionCode) {
	switch {
	case arg == "1":  // print compressed otuput
		act.Print = PrintCompact
	case arg == "2":  // print full otuput
		act.Print = PrintFull
	case arg == "a":  // ascending order
		act.Sort = SortAsc
	case arg == "b":  // browse (print only, no select)
		act.Main = MainActView
		act.Sub = SubActView
	case arg == "d":  // descending order
		act.Sort = SortDesc
	case arg == "e":  // select, edit
		act.Main = MainActView
		act.Sub = SubActEdit
	case arg == "h":  // help
		act.Main = MainActHelp
	case arg == "i":  // init
		act.Main = MainActInit
	case arg == "l":  // match loose
		act.Match = MatchLoose
	case arg == "m":  // Demo.  #TODO
		act.Main = MainActDemo
	case arg == "n":  // create entry
		act.Main = MainActCreate
	case arg == "p":  // select, pipe
		act.Main = MainActView
		act.Sub = SubActPipe
	case arg == "s":  // match strict
		act.Match = MatchStrict
	// case arg == "t":  // view tags
	// 	act = []int{4, 2, 0, 0, 0}
	// case arg == "v":  // view values
	// 	act = []int{4, 1, 0, 0, 0}
	case arg == "x":  // select, delete
		act.Main = MainActView
		act.Sub = SubActDelete
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized short-form option `%v`", arg)
	}
}

// updateActionCodeFromWord is just like `updateActionCodeFromChar`
// except it acts on long-form options.
func updateActionCodeFromWord(arg string, act *ActionCode) {
	switch {
	case arg == "asc":
		act.Sort = SortAsc
	case arg == "browse":
		act.Main = MainActView
		act.Sub = SubActView
	case arg == "desc":
		act.Sort = SortDesc
	case arg == "delete":
		act.Main = MainActView
		act.Sub = SubActDelete
	case arg == "demo":  // Demo.  #TODO
		act.Main = MainActDemo
	case arg == "edit":
		act.Main = MainActView
		act.Sub = SubActEdit
	case arg == "help":
		act.Main = MainActHelp
	case arg == "init":
		act.Main = MainActInit
	case arg == "loose":
		act.Match = MatchLoose
	case arg == "new":
		act.Main = MainActCreate
	case arg == "one-line":
		act.Print = PrintCompact
	case arg == "pipe":
		act.Main = MainActView
		act.Sub = SubActPipe
	case arg == "strict":
		act.Match = MatchStrict
	// case arg == "tags":
	// 	act = []int{4, 2, 0, 0, 0}
	case arg == "two-line":
		act.Print = PrintFull
	// case arg == "vals":
	// 	act = []int{4, 1, 0, 0, 0}
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized long-form option `%v`", arg)
	}
}
