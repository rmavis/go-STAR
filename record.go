package main

import (
    "bufio"
	"fmt"
	"os"
	"sort"
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



func readRecords(args []string, match_lim int) []Record {
	var matcher func(Record) float64
	var sorts func([]Record)

	if len(args) == 0 {
		matcher = func(_ Record) float64 {return 1.0}
		sorts = func(records []Record) {
			sort.Sort(sort.Reverse(ByDateCreated(records)))
		}
	} else {
		matcher = makeMatcher(args, match_lim)
		sorts = func(records []Record) {
			sort.Sort(sort.Reverse(ByMatchRate(records)))
		}
	}

	file_name := DefaultStoreFilePath()

	records := readRecordsFromStore(file_name, matcher)
	sorts(records)

	return records

	// printRecords(records)
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



// The idea for the matcher function is to check the value and the
// tags for each argument. For each match, a match rate will be
// added to a collection. If the number of matches is greater than
// the given limit, then the aggregate of the rates will be returned
// as the record's overall match rate. The aggregate is used instead
// of the average (sum of match rates / number of matches) because
// it makes sense that a record that matches multiple times should
// rate higher than those that don't.
func makeMatcher(terms []string, lim int) func(Record) float64 {
	matcher := func(record Record) float64 {
		var match_rates []float64
		matches := 0

		strs := make([]string, 1, len(record.Tags) + 1)
		strs[0] = record.Value
		for o := 0; o < len(record.Tags); o++ {
			strs = append(strs, record.Tags[o])
		}
		str_agg := strings.Join(strs, string(UnitSeparator))
		// fmt.Printf("Aggregate line: %v\n", str_agg)

		for o := 0; o < len(terms); o++ {
			mult := strings.Count(str_agg, terms[o])

			if mult == 0 {
				match_rates = append(match_rates, 0.0)
			} else {
				matches += 1
				match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) * float64(mult)) / float64(len([]rune(str_agg)))))
			}

			// fmt.Printf("Checking %v\n", record)

			// if strings.Contains(record.Value, terms[o]) {
			// 	// fmt.Printf("V %v contains %v\n", record.Value, terms[o])
			// 	matches += 1
			// 	match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) / float64(len([]rune(record.Value)))) * float64(strings.Count(record.Value, terms[o]))))
			// } else {
			// 	match_rates = append(match_rates, 0.0)
			// }

			// for i := 0; i < len(record.Tags); i++ {
			// 	if strings.Contains(record.Tags[i], terms[o]) {
			// 		matches += 1
			// 		// fmt.Printf("T %v contains %v\n", record.Tags[i], terms[o])
			// 		match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) / float64(len([]rune(record.Tags[i])))) * float64(strings.Count(record.Tags[i], terms[o]))))
			// 	} else {
			// 		match_rates = append(match_rates, 0.0)
			// 	}
			// }
		}

		if matches < lim {
			return 0.0
		} else {
			match_rate := 0.0

			for o := 0; o < len(match_rates); o++ {
				match_rate += match_rates[o]
			}

			return (match_rate / float64(len(match_rates)))
			// return match_rate
		}
	}

	return matcher
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
	// Print a prompt
	// Read user input
	// Split on spaces and commas, filter for integers, sort, make unique
	// Return resulting slice
	return []int{1, 2}
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



func makeRecordReviewer(act func([]Record)) func([]Record) {
	reviewer := func(records []Record) {
		printRecords(records)

		input := promptForWantedRecord()
		wanted := getWantedRecords(records, input)

		if len(wanted) == 0 {
			noRecordsWanted()
		} else {
			act(wanted)
		}
	}

	return reviewer
}



func makeRecordPiper(act string) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act)
	}

	return piper
}



func pipeRecordsToPbcopy(records []Record) {
	pipeRecordsToExternalTool(records, PbcopyPath);
}

func pipeRecordsToOpen(records []Record) {
	pipeRecordsToExternalTool(records, OpenPath);
}

func pipeRecordsToExternalTool(records []Record, tool string) {  // #TODO
	for _, r := range records {
		fmt.Printf("Would pipe record value (%v) to tool (%v)\n", r.Value, tool)
	}
}



func editRecords(records []Record) {
}



func deleteRecords(records []Record) {
}



func noRecordsWanted() {
	fmt.Printf("Nothing wanted.\n")
}
