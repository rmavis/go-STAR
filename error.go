package main

import (
	"fmt"
	"os"
)


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

// printInternalActionCodeError receives an ActionCode and prints
// an error message.
func printInternalActionCodeError(act *ActionCode) {
	fmt.Fprintf(os.Stderr, "There's a problem with the action code (%v).\n", act)
}
