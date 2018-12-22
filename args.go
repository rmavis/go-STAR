package main

import (
	"strings"
)



// ArgPair is a structure that pairs a string -- the arg given by the
// user on the command line -- with a function that will parse that
// arg and return an appropriate action code.
type ArgPair struct {
	Arg string
	Parse func (string) []int
}


// actionCodeHelpCommands returns the action code for printing the
// basic help message. This is necessary because Go doesn't allow
// slices to be constants.
func actionCodeHelpCommands() []int {
	return []int{3, 1, 0}
}


// actionCodeDefaultSearch returns the action code for the basic
// search functionality. This is necessary because Go doesn't allow
// slices to be constants.
func actionCodeDefaultSearch() []int {
	return []int{1, 0, 0}
}





// parseArgs takes a slice of strings, being the command line args,
// and returns a slice of ints that, together, encode the user's
// intent, and a slice of strings that, if present, will affect the
// action -- as search terms, or new inputs, or whatever.
func parseArgs(args []string) ([]int, []string) {
	// If there are no args, then assume help is needed.
	if len(args) == 0 {
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

	// This is the default action, set here in case the only args are
	// search terms, which is likely the most common use case. Other
	// arguments can modify this. If they do, the final action will
	// be the merged result of the last with all usable aspects of
	// prior actions.
	act := actionCodeDefaultSearch()

	// This will contain the strings to search for, AKA non-action
	// arguments, if any are given.
	var strs []string

out:
	for o := 0; o < len(args); o++ {
		// Command switches start with dashes.
		if args[o][0] == '-' {
			arg := strings.ToLower(args[o])

			// Long-form switches start with two.
			if arg[1] == '-' {
				arg = strings.Replace(arg, string('-'), "", -1)

				if _, present := args_ref[arg]; present == false {
					args_ref[arg] = true
					args_use = append(args_use, ArgPair{arg, getActFromWord})
				}
			} else {
				// Short-form switches start with one.
				arg = strings.Replace(arg, string('-'), "", -1)

				for i := 0; i < len(arg); i++ {
					char := string(arg[i])

					if _, present := args_ref[char]; present == false {
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
// represents an argument and returns a slice of ints that encodes
// the corresponding action.
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



// getActFromChar receives a string, being a short-form command-line
// switch, and returns a slice of ints that encodes the corresponding
// action. 0s in the action code indicate that an int from a prior
// action code can be merged in at that location.
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



// getActFromWord receives a string, being a long-form command-line
// switch, and returns a slice of ints that encodes the corresponding
// action. 0s in the action code indicate that an int from a prior
// action code can be merged in at that location.
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
	case arg == "vals":
		act = []int{4, 1, 0}
	default:
		act = []int{0, 0, 0}
	}

	return act
}



// mergeActionCodes receives an action code and a list of action
// codes and returns an action code. For each 0 in the single action
// code, a non-zero in that place will be looked for in each code in
// the list, and the first non-zero will be substituted in its place.
// A non-zero int in an action code indicates a stated intent.
// The list should be considered a pool that can be merged into the
// single. It should be ordered oldest to most recent.
func mergeActionCodes(act []int, acts [][]int) []int {
	new := act

	for o := 0; o < len(new); o++ {
		if new[o] == 0 {
		out:
			for i := len(acts) - 1; i >= 0; i-- {
				// The zeroth position is checked because only codes
				// for the same action can be merged.
				if acts[i][0] == new[0] && acts[i][o] != 0 {
					new[o] = acts[i][o]
					break out;
				}
			}
		}
	}

	return new
}
