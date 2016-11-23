package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)





// makeRecordSelector receives a prompt verb and a record-action
// function and returns an action function that prints the prompts the user
// for the .
func makeRecordSelector(prompt_verb string, act func([]Record)) func([]Record) {
	selector := func(records []Record) {
		printRecordsToStdout(records)

		input := promptForWantedRecord(prompt_verb)
		wanted := getWantedRecords(records, input)

		if len(wanted) == 0 {
			willDoNothing(prompt_verb)
		} else {
			act(wanted)
		}
	}

	return selector
}



// promptForWantedRecord prints a prompt containing the given verb to
// stdout. It will collect the user's input, remove whitespace, and
// return it.
func promptForWantedRecord(verb string) string {
	fmt.Printf("%v%v these records: ", strings.ToUpper(string(verb[0])), string(verb[1:]))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input)
}



// getWantedRecords receives a slice of Records and a string. It will
// check that string and return either a sub-slice of the Records
// indicated by the numbers in the string, or the entire slice if the
// string reads "all".
func getWantedRecords(records []Record, input string) []Record {
	if strings.ToLower(input) == "all" {
		return records
	} else {
		ints := getIntsFromInput(input)

		var wanted []Record
		max := len(records)

		for _, i := range ints {
			if i <= max {
				wanted = append(wanted, records[(i - 1)])
			}
		}

		return wanted
	}
}



// getIntsFromInput receives a string and returns a slice of the
// unique ints that string contains. Ints in the string can be
// separated by spaces and/or commas.
func getIntsFromInput(input string) []int {
	var clean []int
	ref := make(map[string]bool)

	// Note that the comma is in double-quotes. Those make it a
	// string. Single-quotes make it a rune, which can't be used as
	// a parameter to `Replace`.
	parts := strings.Split(strings.Replace(input, ",", " ", -1), " ")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)

		if num, err := strconv.Atoi(trimmed); err == nil {
			if _, in_ref := ref[trimmed]; in_ref == false {
				clean = append(clean, num)
				ref[trimmed] = true
			}
		}
	}

	return clean
}



// willDoNothing prints a message containing the given verb to stdout
// indicating to the user that no action will be taken.
func willDoNothing(verb string) {
	fmt.Printf("Will %v nothing.\n", verb)
}
