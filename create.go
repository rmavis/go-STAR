package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)





// makeCreateAction returns the Create action function. It requires
// the user's config and the terms given on the command line.
func makeCreateAction(conf Config, terms []string) func() {
	action := func() {
		if len(terms) == 0 {
			printTermCreationError()
		} else {
			record := makeRecordFromInput(terms)
			appendRecordsToFile(conf.Store, []Record{record})
		}
	}

	return action
}



// makeRecordFromInput transforms the given terms, adds initial meta-
// data, and returns a fully-formed Record.
func makeRecordFromInput(terms []string) Record {
	val := terms[0]
	tags := terms[1:]

	s := strconv.FormatInt(time.Now().Unix(), 10)
	times := []string{s, "0", "0"}

	return Record{val, tags, times, 0.0}
}



// appendRecordsToFile appends the given records, one by one, to the
// file named by the given string.
func appendRecordsToFile(file_name string, records []Record) {
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY, 0600)
	checkForError(err)
	defer file.Close()

	for _, record := range records {
		saveRecordToFile(file, record)
	}
}



func printTermCreationError() {
	fmt.Println("A new entry needs a value and any number of tags. Example:\n  $ star -n value tag1 tag2 tag3")
}
