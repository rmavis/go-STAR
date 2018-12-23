package main

import (
	"fmt"
	"os"
	"strings"
)


// actionCodeDefaultView returns the action code for the basic
// search functionality. This is necessary because Go doesn't allow
// slices to be constants.
func actionCodeDefaultView() []int {
	return []int{1, 0, 0, 0, 0}
}

// parseArgs takes a slice of strings, being the command line args,
// and returns a slice of ints that, together, encode the user's
// intent, and a slice of strings that, if present, will affect the
// action -- as search terms, or new inputs, or whatever. The slice
// returned should be interpreted like:
// 0 = major action (select, create, etc)
// 1 = sub-action (edit matches, delete matches, etc)
// 3 = match mode (loose, strict)
// 4 = sort order (high to low, low to high)
// 5 = output format (compressed, full)
func parseArgs(args []string) ([]int, []string) {
	act := actionCodeDefaultView()
	var strs []string

out:
	for o := 0; o < len(args); o++ {
		// Command options start with dashes.
		if args[o][0] == '-' {
			arg := strings.ToLower(args[o])

			if arg[1] == '-' {  // Long-form options start with two.
				arg = strings.Replace(arg, string('-'), "", -1)
				act = mergeActionCodes(act, getActFromWord(arg))
			} else {  // Short-form options start with one.
				arg = strings.Replace(arg, string('-'), "", -1)
				for i := 0; i < len(arg); i++ {
					char := string(arg[i])
					act = mergeActionCodes(act, getActFromChar(char))
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

// getActFromChar receives a string, being a short-form command-line
// option, and returns a slice of ints that encodes the corresponding
// action. 0s in the action code indicate that an int from a prior
// action code can be merged in at that location.
func getActFromChar(arg string) []int {
	var act []int

	switch {
	case arg == "1":  // print compressed otuput
		act = []int{0, 0, 0, 0, 1}
	case arg == "2":  // print full otuput
		act = []int{0, 0, 0, 0, 2}
	case arg == "a":  // ascending order
		act = []int{0, 0, 0, 2, 0}
	case arg == "b":  // browse (print only, no select)
		act = []int{1, 1, 0, 0, 0}
	case arg == "d":  // descending order
		act = []int{0, 0, 0, 1, 0}
	case arg == "e":  // select, edit
		act = []int{1, 3, 0, 0, 0}
	case arg == "h":  // help
		act = []int{3, 0, 0, 0, 0}
	case arg == "i":  // init
		act = []int{5, 0, 0, 0, 0}
	case arg == "l":  // match loose
		act = []int{0, 0, 1, 0, 0}
	case arg == "m":  // Demo.  #TODO
		act = []int{6, 0, 0, 0, 0}
	case arg == "n":  // create entry
		act = []int{2, 0, 0, 0, 0}
	case arg == "p":  // select, pipe
		act = []int{1, 2, 0, 0, 0}
	case arg == "s":  // match strict
		act = []int{0, 0, 2, 0, 0}
	// case arg == "t":  // view tags
	// 	act = []int{4, 2, 0, 0, 0}
	// case arg == "v":  // view values
	// 	act = []int{4, 1, 0, 0, 0}
	case arg == "x":  // select, delete
		act = []int{1, 4, 0, 0, 0}
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized option `%v`", arg)
		act = []int{0, 0, 0, 0, 0}
	}

	return act
}

// getActFromWord receives a string, being a long-form command-line
// option, and returns a slice of ints that encodes the corresponding
// action. 0s in the action code indicate that an int from a prior
// action code can be merged in at that location.
func getActFromWord(arg string) []int {
	var act []int

	switch {
	case arg == "asc":
		act = []int{0, 0, 0, 2, 0}
	case arg == "browse":
		act = []int{1, 1, 0, 0, 0}
	case arg == "desc":
		act = []int{0, 0, 0, 1, 0}
	case arg == "delete":
		act = []int{1, 4, 0, 0, 0}
	case arg == "demo":  // Demo.  #TODO
		act = []int{6, 0, 0, 0, 0}
	case arg == "edit":
		act = []int{1, 3, 0, 0, 0}
	case arg == "help":
		act = []int{3, 0, 0, 0, 0}
	case arg == "init":
		act = []int{5, 0, 0, 0, 0}
	case arg == "loose":
		act = []int{0, 0, 1, 0, 0}
	case arg == "new":
		act = []int{2, 0, 0, 0, 0}
	case arg == "one-line":
		act = []int{0, 0, 0, 0, 1}
	case arg == "pipe":
		act = []int{1, 2, 0, 0, 0}
	case arg == "strict":
		act = []int{0, 0, 2, 0, 0}
	// case arg == "tags":
	// 	act = []int{4, 2, 0, 0, 0}
	case arg == "two-line":
		act = []int{0, 0, 0, 0, 2}
	// case arg == "vals":
	// 	act = []int{4, 1, 0, 0, 0}
	default:
		fmt.Fprintf(os.Stderr, "Unrecognized option `%v`", arg)
		act = []int{0, 0, 0, 0, 0}
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
func mergeActionCodes(act_a []int, act_b []int) []int {
	//fmt.Fprintf(os.Stderr, "MERGING ACTION CODES: `%v` into `%v`", act_b, act_a)
	new := act_a

	for o := 0; o < len(act_b); o++ {
		if act_b[o] > 0 {
			new[o] = act_b[o]
		}
	}

	//fmt.Fprintf(os.Stderr, "MERGED ACTION CODES: `%v`", new)
	return new
}
