package main

import (
    "os"
)





func main() {
	args := os.Args[1:]
	action := makeAction(args)
	action()
}



func checkForError(e error) {
	if e!= nil {
		panic(e)
	}
}
