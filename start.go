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
		action = func() {fmt.Printf("Would make `help` action.")}
	case action_code[0] == 4:  // dump
		if action_code[1] == 1 {
			action = makeSearchAction(readConfig(), action_code, terms)
		} else {
			action = func() {fmt.Printf("Would make `dump` action.")}  // #TODO
		}
	case action_code[0] == 5:  // initialize
		action = func() {fmt.Printf("Would make `init` action.")}
	case action_code[0] == 6:  // demo
		action = func() {fmt.Printf("Would make `demo` action.")}
	default:
		action = func() {fmt.Printf("Would make `error` action.")}
	}

	action()
}



// checkForError is a convenience function to cope with Go's idiom
// of returning an error message if a function call fails. So rather
// than doing this all the time:
//   old_count, err := strconv.Atoi(record.Meta[2])
//   if e!= nil {
//     panic(e)
//   }
// there can just be this:
//   old_count, err := strconv.Atoi(record.Meta[2])
//   checkForError(err)
func checkForError(e error) {
	if e!= nil {
		panic(e)
	}
}
