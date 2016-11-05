package main

import (
	"fmt"
	"strings"
)





type ArgPair struct {
	Arg string
	Parse func (string) []int
}



func actionCodeHelpCommands() []int {
	return []int{3, 1, 0}
}


func actionCodeDefaultSearch() []int {
	return []int{1, 0, 0}
}





// parseArgs takes a slice of strings, being the command line args,
// and returns a slice of ints that, together, encode the user's
// intent, and a slice of strings that, if present, will affect the
// action -- as search terms, or new inputs, or whatever.
func parseArgs(args []string) ([]int, []string) {
	if len(args) == 0 {
		fmt.Printf("No args given: using code %v\n", actionCodeHelpCommands())
		var empty_str []string
		return actionCodeHelpCommands(), empty_str
	}

	// The keys of this will be unique action arguments. It's used
	// as a reference -- actions are only found for unique arguments,
	// but the order of the arguments also matters.
	args_ref := make(map[string]bool)

	// So this will contain each (unique) argument in the order in
	// which it first occurs. For example, the arguments
	//   -acacaaacca --all --all --commands --all
	// will be captured as
	//   ['a', 'c', 'all', 'commands']
	var args_use []ArgPair

	// This is the default action, set here in case the only args
	// are search terms. If arguments modify this, the result will
	// be a sort of merged-up action combining the last action with
	// all usable aspects of prior actions.
	act := actionCodeDefaultSearch()

	// This will contain the strings to search for, AKA non-action
	// arguments, if any are given.
	var strs []string

out:
	for o := 0; o < len(args); o++ {
		// Arguments start with dashes.
		if args[o][0] == '-' {
			arg := strings.ToLower(args[o])

			// Long-form args start with two dashes.
			if arg[1] == '-' {
				arg = strings.Replace(arg, string('-'), "", -1)
				_, present := args_ref[arg]

				if present == false {
					args_ref[arg] = true
					args_use = append(args_use, ArgPair{arg, getActFromWord})
				}
			} else {
				// Short-form args start with one.
				arg = strings.Replace(arg, string('-'), "", -1)

				for i := 0; i < len(arg); i++ {
					char := string(arg[i])
					_, present := args_ref[char]

					if present == false {
						args_ref[char] = true
						args_use = append(args_use, ArgPair{char, getActFromChar})
					}
				}
			}
		} else {
			strs = args[o:]
			break out;
		}
	}

	if len(args_use) > 0 {
		var acts [][]int

		for o := 0; o < len(args_use); o++ {
			act = args_use[o].Parse(args_use[o].Arg)
			acts = append(acts, act)

			if o > 0 {
				act = mergeActionCodes(act, acts)
			}
		}
	}

	return act, strs
}



// checkShortFormArgs receives a string in which each character
// represents an argument.
func checkShortFormArgs(args string) []int {
	if len(args) == 1 {
		return getActFromChar(string(args[0]))
	} else {
		var act []int
		var acts [][]int

		for o := 0; o < len(args); o++ {
			act = getActFromChar(string(args[o]))
			acts = append(acts, act)
			act = mergeActionCodes(act, acts)
		}

		return act
	}
}



func getActFromChar(arg string) []int {
	var act []int

	switch {
	case arg == "a":  // dump vals (all)
		act = []int{4, 1, 0}
	case arg == "b":  // search, print only (browse)
		act = []int{1, 5, 0}
	case arg == "c":  // search, copy
		act = []int{1, 1, 0}
	case arg == "d":  // search, delete
		act = []int{1, 4, 0}
	case arg == "e":  // search, edit
		act = []int{1, 3, 0}
	case arg == "h":  // help, commands
		act = []int{3, 1, 0}
	case arg == "i":  // init
		act = []int{5, 0, 0}
	case arg == "l":  // search, match loose
		act = []int{1, 0, 1}
	case arg == "m":  // This will result in a lame demo.  #TODO
		act = []int{6, 0, 0}
	case arg == "n":  // create entry
		act = []int{2, 0, 0}
	case arg == "o":  // search, open
		act = []int{1, 2, 0}
	case arg == "r":  // help, readme
		act = []int{3, 2, 0}
	case arg == "s":  // search, match strict
		act = []int{1, 0, 2}
	case arg == "t":  // dump tags
		act = []int{4, 2, 0}
	case arg == "x":  // help, customization
		act = []int{3, 3, 0}
	default:
		act = []int{0, 0, 0}
	}

	return act
}



func getActFromWord(arg string) []int {
	var act []int

	switch {
	case arg == "all":
		act = []int{4, 1, 0}
	case arg == "browse":
		act = []int{1, 5, 0}
	case arg == "commands":
		act = []int{3, 1, 0}
	case arg == "customization":
		act = []int{3, 3, 0}
	case arg == "delete":
		act = []int{1, 4, 0}
	case arg == "edit":
		act = []int{1, 3, 0}
	case arg == "flags":
		act = []int{3, 1, 0}
	case arg == "help":
		act = []int{3, 1, 0}
	case arg == "init":
		act = []int{5, 0, 0}
	case arg == "loose":
		act = []int{1, 0, 1}
	case arg == "new" || arg == "create":
		act = []int{2, 0, 0}
	case arg == "open":
		act = []int{1, 2, 0}
	case arg == "readme":
		act = []int{3, 2, 0}
	case arg == "strict":
		act = []int{1, 0, 2}
	case arg == "tags":
		act = []int{4, 2, 0}
	default:
		act = []int{0, 0, 0}
	}

	return act
}



// A non-zero int in an action code indicates a stated intent.
// A zero indicates that the previous stated intent should be used
// in that place or, if there are none, the appropriate default.
func mergeActionCodes(act []int, acts [][]int) []int {
	new := act

	// fmt.Printf("Merging %v into %v\n", act, acts)

	for o := 0; o < len(new); o++ {
		if new[o] == 0 {
		out:
			for i := len(acts) - 1; i >= 0; i-- {
				if acts[i][0] == new[0] && acts[i][o] != 0 {
					new[o] = acts[i][o]
					break out;
				}
			}
		}
	}

	return new
}
