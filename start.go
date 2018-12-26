package main

import (
	"fmt"
	"os"
)


// main runs the show.
func main() {
	var action func()

	act, terms := parseArgs(os.Args[1:])
	switch {
	case act.Main == MainActView:
		action = makeSearchAction(readConfig(), act, terms)
	case act.Main == MainActCreate:
		action = makeCreateAction(readConfig(), terms)
	case act.Main == MainActHelp:
		action = func() {printUsageInformation()}
	case act.Main == MainActInit:
		action = makeInitializer(terms)
	case act.Main == MainActDemo:
		action = func() {fmt.Printf("Would make `demo` action.")}  // #TODO
	default:
		action = func() {printInternalActionCodeError(act)}
	}

	action()
}
