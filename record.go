package main

import (
    "bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	// "unicode/utf8"
)


type Record struct {
	Value string
	Tags []string
	Meta []string
	MatchRate float64
}



// These are the ASCII Group, Record, and Unit separator characters.
// They are of type `rune`.
const GroupSeparator = ''   // Separates records
const RecordSeparator = ''  // Separates parts of records (the value from the tags)
const UnitSeparator = ''    // Separates parts of parts (each tag from each tag)





func readRecordsFromStore(file_name string, getMatchInfo func(Record) (float64, bool)) []Record {
	var records []Record

	act := func(record Record) {
		match_rate, matches := getMatchInfo(record)

		if matches {
			record.MatchRate = match_rate
			records = append(records, record)
			// fmt.Printf("Record matches: %v / %v / %v / %v\n", record.Value, record.Tags, record.Meta, record.MatchRate)
		}
	}

	forEachRecordInStore(file_name, act)

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



func makeRecordFromParts(entry []string) Record {
	return Record{entry[0], splitField(entry[1]), splitField(entry[2]), 0.0}
}



func splitEntry(entry string) []string {
	fields := strings.Split(strings.TrimSuffix(entry, string(GroupSeparator)), string(RecordSeparator))
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



func joinRecord(record Record) string {
	parts := []string{
		record.Value,
		string(RecordSeparator),
		strings.Join(record.Tags, string(UnitSeparator)),
		string(RecordSeparator),
		strings.Join(record.Meta, string(UnitSeparator)),
		string(GroupSeparator),
		"\n",
	}

	return strings.Join(parts, "")
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
		// fmt.Printf("%v%v) %v%v\n%v%v\n",
		// 	spaces_top, (o + 1), records[o].MatchRate, records[o].Value,
		// 	spaces_bot, strings.Join(records[o].Tags, ", "))
	}
}



func promptForWantedRecord(verb string) []int {
	fmt.Printf("%v%v these records: ", strings.ToUpper(string(verb[0])), string(verb[1:]))
	// fmt.Print("Enter the number(s) of the record(s) you want: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	ints := cleanWantedRecordsInput(input)
	// fmt.Printf("Got (%v)\n", ints)

	return ints
}



func cleanWantedRecordsInput(input string) []int {
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



func makeRecordSelector(prompt_verb string, act func([]Record)) func([]Record) {
	selector := func(records []Record) {
		printRecords(records)

		input := promptForWantedRecord(prompt_verb)
		wanted := getWantedRecords(records, input)

		if len(wanted) == 0 {
			noRecordsWanted(prompt_verb)
		} else {
			act(wanted)
		}
	}

	return selector
}



func makeActAndUpdater(conf Config, act func([]Record)) func([]Record) {
	saver := func(file *os.File, record Record) {
		file.WriteString(joinRecord(record))
	}

	updater := func(records []Record) {
		updateWantedRecords(records)
		updateStoreFile(conf.Store, records, saver)
		act(records)
	}

	return updater
}



func makeRecordPiper(act string) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act)
	}

	return piper
}



func makeDeleter(conf Config) func([]Record) {
	dels := func(_ *os.File, _ Record) {
		// fmt.Printf("Not including entry in backup (%v)\n", record)
	}

	deleter := func(records []Record) {
		updateStoreFile(conf.Store, records, dels)
	}

	return deleter
}



func editRecords(records []Record) {
	// Make tmp file
	// Print instructions to tmp file
	// Print records to tmp file
	// Open editor in tmp file
	// Wait for editor to close
	// Open temp file
	// For each line, check for patterns of entry or tags
	//   Numbers in each entry line are the indices for the matching records
	//   Those numbers matter
	// Make a slice of new Records
	//   Update the values and tags for each index
	//   Add new records for indices beyond the max
	//   Remove records for indices no longer specified
	// Compare new slice with existing?
	//   To make "wanted" slice of records?
	// Update store as normal with modified matching records as "wanted" records
}



func updateWantedRecords(records []Record) {
	now := strconv.FormatInt(time.Now().Unix(), 10)

	for _, record := range records {
		switch {
		case len(record.Meta) == 3:
			old_count, err := strconv.Atoi(record.Meta[2])
			checkForError(err)
			record.Meta[2] = strconv.Itoa(old_count + 1)
			record.Meta[1] = now
		default:
			fmt.Printf("Record is missing metadata! (%v) (%v)\n", record.Value, record.Meta)
		}
		// if utf8.RuneCountInString(record.Meta)
	}
}



func noRecordsWanted(verb string) {
	fmt.Printf("Will %v nothing.\n", verb)
}
