package main

import (
	"fmt"
	"os"
)


// main runs the show.
func main() {
	// This gets the ActionCode and terms from the command line
	// arguments, if any exist.
	act, terms := parseArgs(os.Args[1:])

	// `action` will be a function that, like `main`, will be called
	// with no arguments and that will return nothing.
	var action func()

	switch {
	case act.Main == MainActView:  // view
		action = makeSearchAction(readConfig(), act, terms)
	case act.Main == MainActCreate:  // create
		action = makeCreateAction(readConfig(), terms)
	case act.Main == MainActHelp:  // help
		action = func() {printUsageInformation()}
	case act.Main == MainActInit:  // initialize
		action = makeInitializer(terms)
	case act.Main == MainActDemo:  // demo
		action = func() {fmt.Printf("Would make `demo` action.")}  // #TODO
	default:
		action = func() {printInternalActionCodeError(act)}
	}

	action()
}
