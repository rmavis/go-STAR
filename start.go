package main

import (
	"fmt"
	"os"
)


// main runs the show.
func main() {
	// This gets the ActionCode and terms from the command line
	// arguments, if any exist.
	action_code, terms := parseArgs(os.Args[1:])

	// `action` will be a function that, like `main`, will be called
	// with no arguments and that will return nothing.
	var action func()

	switch {
	case action_code.Main == MainActView:  // view
		action = makeSearchAction(readConfig(), action_code, terms)
	case action_code.Main == MainActCreate:  // create
		action = makeCreateAction(readConfig(), terms)
	case action_code.Main == MainActHelp:  // help
		action = func() {printUsageInformation()}
	case action_code.Main == MainActInit:  // initialize
		action = makeInitializer(terms)
	case action_code.Main == MainActDemo:  // demo
		action = func() {fmt.Printf("Would make `demo` action.")}  // #TODO
	default:
		action = func() {printInternalActionCodeError(action_code)}
	}

	action()
}
