package main

import (
	"fmt"
	"strings"
)


/*

The functions in this file will build and pass around a slice that
contains three ints that, together, represent the user's intention
as specified in the command line arguments.

The first int is a code indicating the command/operation:
1: search
2: create
3: help
4: dump
5: initialize
6: demo

If the first int is a verb, the second int is an adverb. Only some
verbs can be modified:
1 (search):
  0: default (read from config)
  1: pbcopy
  2: open
  3: edit
  4: delete
3 (help):
  1: commands
  2: readme
  3: customization
  4: examples
4 (dump):
  1: values
  2: tags

And if the second int is an adverb, the third is an adjective:
1 (search):
  0: default (read from config)
  1: loose
  2: strict

A 0 in any place means that no relevant argument was given and that
a sensible default should be used, if applicable.

*/




func makeAction(args []string) func() {
	action_code, terms := parseArgs(args)
	conf := readConfig()

	var action func()
	switch {
	case action_code[0] == 1:  // search
		action = makeSearchAction(conf, action_code, terms)

	case action_code[0] == 2:  // create
		action = func() {fmt.Printf("Would make `create` action.")}
	case action_code[0] == 3:  // help
		action = func() {fmt.Printf("Would make `help` action.")}
	case action_code[0] == 4:  // dump
		if action_code[1] == 1 {
			action = makeSearchAction(conf, action_code, terms)
		} else {
			action = func() {fmt.Printf("Would make `dump` action.")}
			// #TODO
		}

	case action_code[0] == 5:  // initialize
		action = func() {fmt.Printf("Would make `init` action.")}
	case action_code[0] == 6:  // demo
		action = func() {fmt.Printf("Would make `demo` action.")}
	default:
		action = func() {fmt.Printf("Would make `error` action.")}
	}

	return action
}



func makeSearchAction(conf Config, action_code []int, terms []string) func() {
	var match_act func ([]Record)
	var match_lim int

	// The 1st position indicates the action to take on the matches.
	switch {
	case action_code[1] == 0:  // Read from config.
		switch {
		case conf.Action == "copy":  // "copy" and "open" are shortcuts.
			match_act = makeRecordReviewer(pipeRecordsToPbcopy)
		case conf.Action == "open":
			match_act = makeRecordReviewer(pipeRecordsToOpen)
		case conf.Action == "":      // If nothing, then just show them.
			match_act = printRecords
		default:                     // Any external command can be specified.
			piper := makeRecordPiper(conf.Action)
			match_act = makeRecordReviewer(piper)
		}
	case action_code[1] == 1:  // pbcopy
		match_act = makeRecordReviewer(pipeRecordsToPbcopy)
	case action_code[1] == 2:  // open
		match_act = makeRecordReviewer(pipeRecordsToOpen)
	case action_code[1] == 3:  // edit
		match_act = editRecords
	case action_code[1] == 4:  // delete
		match_act = deleteRecords
	default:  // Bork.
		match_act = printRecords
	}

	// The 2nd position indicates the match mode.
	switch {
	case action_code[0] == 4:  // All.
		match_lim = 0
	case action_code[2] == 0:  // Read from config.
		if conf.FilterMode == "loose" {
			match_lim = 1
		} else {
			match_lim = len(terms)
		}
	case action_code[2] == 1:  // loose
		match_lim = 1
	case action_code[2] == 2:  // strict
		match_lim = len(terms)
	default:  // Bork.
		fmt.Printf("Invalid match code (%v). Using '1'.\n", action_code[2])
		match_lim = 1
	}

	action := func() {
		records := readRecords(terms, match_lim)
		match_act(records)
	}

	return action
}





//////////////////////////////





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
	case arg == "n":  // new entry
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
	case arg == "new":
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




/*

ARGUMENT TYPES

- Matching
  l (match loosely -- or)
  s (match strictly -- and)

- Action
  o (pipe value to `open`)
  c (pipe value to `pbcopy`)
  (((would be cool if user could pipe value to arbitrary program)))
  d (delete)
  e (edit)

- Static messages
  c (help [commands])
  f (help [flags])
  h (help)
  r (readme)
  xh (extra help)
  xr (extra readme)
  x (examples)

- Dynamic messages (dumps?)
  a (all)
  t (tags)

- Creation
  i (init -- `touch` store file)
  n (add new entry)

- Compound
  m (demo)



USAGE FORMS

- Find and operate on exiting entries:
  star -[ls][ocde] tag(s)
  Operation (Verb): search.
  Match mode (Adverb): determined by [ls] or default.
  Action: determined by [ocde] or default.

- See a help message
  star -[cfh][r][x]
  Operation (Verb): help message.
  Variation: determined by optional presence of `x`.

- See data dumps on existing entries
  star -[at]
  Operation (Verb): dump.
  Dump material: determined by [at].

- Create a new entry
  star -n tag(s) value
  Operation (Verb): create.

- Touch the store file
  star -i
  Operation (Verb): initialize.

- Run the demo
  star -m -[ls][ocde] tag(s)
  Operation (Verb): demo.


FLAGS
  -a, --all        Show all entries.
  -c, --commands,
    -f, --flags,
    -h, --help     Show this message.
  -d, --delete     Delete an entry.
  -e, --edit       Edit an entry.
  -i, --init       Create the ~/.config/star/store file.
  -l, --loose      Match loosely, rather than strictly.
  -m, --demo       Run the demo.
  -n, --new        Add a new entry.
  -o, --open       open the value rather than pbcopy it.
  -p, --copy       pbcopy the value rather than open it.
  -r, --readme,    Show the readme message.
  -s, --strict     Match strictly rather than loosely.
  -t, --tags       Show all tags.
  -x, --examples   Show some examples.
  -xh, -hx,        Show this message with extra details.
  -xr, -rx,        Show the readme message with extra details.

*/
