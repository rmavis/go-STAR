package main

import (
	"fmt"
    "os"
)





func main() {
	args := os.Args[1:]
	action := makeAction(args)
	action()
}



func makeAction(args []string) func() {  // #TODO
	action_code, terms := parseArgs(args)
	conf := readConfig()

	var action func()
	switch {
	case action_code[0] == 1:  // search
		action = makeSearchAction(conf, action_code, terms)
	case action_code[0] == 2:  // create
		action = makeCreateAction(conf, action_code, terms)
		// action = func() {fmt.Printf("Would make `create` action.")}
	case action_code[0] == 3:  // help
		action = func() {fmt.Printf("Would make `help` action.")}
	case action_code[0] == 4:  // dump
		if action_code[1] == 1 {
			action = makeSearchAction(conf, action_code, terms)
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

	return action
}



func checkForError(e error) {
	if e!= nil {
		panic(e)
	}
}
