package main

import (
    "bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


type Record struct {
	Value string
	Tags []string
	Meta []string
	MatchRate float64
}



// These are the ASCII Group, Record, and Unit separator characters.
// They are of type `rune`.
const GroupSeparator = ''
const RecordSeparator = ''
const UnitSeparator = ''



func readRecords(matcher func(Record) float64, sorter func([]Record)) []Record {
	file_name := DefaultStoreFilePath()

	records := readRecordsFromStore(file_name, matcher)
	sorter(records)

	return records
}



func readRecordsFromStore(file_name string, getMatchRate func(Record) float64) []Record {
	file_handle, err := os.Open(file_name)
	checkForError(err)
	defer file_handle.Close()

	reader := bufio.NewReader(file_handle)

	var records []Record

	for {
		entry, last := readNextEntry(reader)
		parts := splitEntry(entry)

		if (doesEntryHaveParts(parts)) {
			record := makeRecord(parts)
			match_rate := getMatchRate(record)

			if match_rate > 0.0 {
				record.MatchRate = match_rate
				records = append(records, record)
				// fmt.Printf("Record matches: %v / %v / %v / %v\n", record.Value, record.Tags, record.Meta, record.MatchRate)
			}
		} else {
			// fmt.Printf("Record is missing components: %v\n", entry)
		}

		if last {
			break;
		}
	}

	return records
}



func readNextEntry(reader *bufio.Reader) (string, bool) {
	record, err := reader.ReadBytes(GroupSeparator)
	last := false

	if err != nil {
		last = true
		// fmt.Printf("Error! %v (%v)\n", err, string(record))
	}

	return strings.TrimSpace(string(record)), last;
}



func makeRecord(entry []string) Record {
	return Record{entry[0], splitField(entry[1]), splitField(entry[2]), 0.0}
}



func splitEntry(entry string) []string {
	fields := strings.Split(entry, string(RecordSeparator))
	return fields
}



func splitField(field string) []string {
	parts := strings.Split(field, string(UnitSeparator))
	return parts
}



func doesEntryHaveParts(entry []string) bool {
	if (len(entry) == 3) {
		return true
	} else {
		return false
	}
}



func printRecords(records []Record) {
	// This is the number of records.
	m := len(records)
	// This is the number of digits in that number.
	n := len(strconv.FormatInt(int64(m), 10))

	// fmt.Printf("Number of records: %v. Number of digits: %v.\n", m, n)

	spaces_bot := strings.Repeat(" ", (n + 2))

	for o := 0; o < m; o++ {
		spaces_top := ""

		if v := (n - len(strconv.FormatInt(int64(o + 1), 10))); v > 0 {
			spaces_top += strings.Repeat(" ", v)
		}

		fmt.Printf("%v%v) %v\n%v%v\n",
			spaces_top, (o + 1), records[o].Value,
			spaces_bot, strings.Join(records[o].Tags, ", "))
	}
}



func promptForWantedRecord() []int {
	fmt.Print("Enter the number(s) of the record(s) you want: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	ints := cleanWantedRecordsInput(input)
	// fmt.Printf("Got (%v)\n", ints)

	// return []int{1, 2}
	return ints
}



func cleanWantedRecordsInput(input string) []int {
	var clean []int
	ref := make(map[string]bool)

	// Note that the comma is in double-quotes. Those make it a
	// string. Single-quotes make it a rune, which can't be used as
	// a parameter to `Replace`.
	parts := strings.Split(strings.Replace(input, ",", "", -1), " ")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)

		if num, err := strconv.Atoi(trimmed); err == nil {
			if _, in_ref := ref[trimmed]; in_ref == false {
				clean = append(clean, num)
				ref[trimmed] = true
			}
		}
	}

	// Should the ints be sorted? Consider it.  #TODO

	return clean
}



func getWantedRecords(records []Record, input []int) []Record {
	var wanted []Record
	max := len(records)

	for _, i := range input {
		if i <= max {
			wanted = append(wanted, records[(i - 1)])
		}
	}

	return wanted
}



func makeRecordSelector(act func([]Record)) func([]Record) {
	selector := func(records []Record) {
		printRecords(records)

		input := promptForWantedRecord()
		wanted := getWantedRecords(records, input)

		if len(wanted) == 0 {
			noRecordsWanted()
		} else {
			act(wanted)
			// Wrapup function (update wanteds' time and count, etc)  #TODO
		}
	}

	return selector
}



func makeRecordPiper(act string) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act)
	}

	return piper
}



func editRecords(records []Record) {
}



func deleteRecords(records []Record) {
}



func noRecordsWanted() {
	fmt.Printf("Nothing wanted.\n")
}
