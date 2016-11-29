package main

import (
	"fmt"
	"os"
)





// main runs the show.
func main() {
	// This gets the action code and terms from the command line
	// arguments, if any exist.
	action_code, terms := parseArgs(os.Args[1:])

	// `action` will be a function that, like `main`, will be called
	// with no arguments and that will return nothing.
	var action func()

	switch {
	case action_code[0] == 1:  // search
		action = makeSearchAction(readConfig(), action_code, terms)
	case action_code[0] == 2:  // create
		action = makeCreateAction(readConfig(), terms)
	case action_code[0] == 3:  // help
		action = func() {printUsageInformation()}
	case action_code[0] == 4:  // dump
		if action_code[1] == 1 {
			action = makeSearchAction(readConfig(), action_code, terms)
		} else {
			action = func() {fmt.Printf("Would make `dump` action.")}  // #TODO
		}
	case action_code[0] == 5:  // initialize
		action = makeInitializer(terms)
	// case action_code[0] == 6:  // demo
	// 	action = func() {fmt.Printf("Would make `demo` action.")}  // #TODO
	default:
		action = func() {printInternalActionCodeError(action_code)}
	}

	action()
}
