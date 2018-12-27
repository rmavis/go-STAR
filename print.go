package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)


func getPrinter(act *ActionCode) func([]Record) {
	if act.Print == PrintValsOnly {
		return printRecordsValuesOnly
	} else if (act.Print == PrintCompact) {
		return printRecordsCompact
	} else {
		return listRecordsToStdout
	}
}

// makeRecordPrintCaller returns a function that will print the
// records it receives to stdout if there are more than zero, else
// it will print a message indicating that there are none.
func makeRecordPrintCaller(printer func([]Record)) func([]Record) {
	caller := func(records []Record) {
		if len(records) == 0 {
			noRecordsMatch()
		} else {
			printer(records)
		}
	}

	return caller
}

// printRecordsFull prints the given slice of Records to the given
// io.Writer in the given format.
func printRecordsFull(out io.Writer, records []Record, format string) {
	// This is the number of records.
	m := len(records)
	// This is the number of digits in that number.
	n := len(strconv.FormatInt(int64(m), 10))
	// This is the number of spaces to print on the bottom line.
	spaces_bot := strings.Repeat(" ", (n + 2))

	for o := 0; o < m; o++ {
		spaces_top := ""

		if v := (n - len(strconv.FormatInt(int64(o + 1), 10))); v > 0 {
			spaces_top += strings.Repeat(" ", v)
		}

		fmt.Fprintf(out, format,
			spaces_top, (o + 1), records[o].Value,
			spaces_bot, strings.Join(records[o].Tags, ", "))
	}
}

// printRecordsCompact receives a slice of Records and writes
// the value of each to stdout.
func printRecordsCompact(records []Record) {
	for o := 0; o < len(records); o++ {
		fmt.Fprintf(os.Stdout, "%v\t{%v}\n", records[o].Value, strings.Join(records[o].Tags, ", "))
	}
}

// printRecordsValuesOnly receives a slice of Records and writes
// the value of each to stdout.
func printRecordsValuesOnly(records []Record) {
	for o := 0; o < len(records); o++ {
		fmt.Fprintf(os.Stdout, "%v\n", records[o].Value)
	}
}


// listRecordsToStdout is a convenience function for printing the
// given records to stdout.
func listRecordsToStdout(records []Record) {
	printRecordsFull(os.Stdout, records, "%v%v) %v\n%v%v\n")
}

// listRecordsToTempFile is a convenience function for printing the
// given records to the given file handle.
func listRecordsToTempFile(records []Record, file *os.File) {
	printRecordsFull(file, records, "%v%v) %v\n%vTags: %v\n\n")
}
