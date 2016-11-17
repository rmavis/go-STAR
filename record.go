package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
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



// func readNextEntry(reader *bufio.Reader) (string, bool) {
// 	record, err := reader.ReadBytes(GroupSeparator)
// 	last := false

// 	if err != nil {
// 		last = true
// 		// fmt.Printf("Error! %v (%v)\n", err, string(record))
// 	}

// 	return strings.TrimSpace(string(record)), last;
// }



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



func removeRecord(records []Record, index int) []Record {
	var new_records []Record

	switch {
	case index == 0:
		new_records = records[1:]
	case index == (len(records) - 1):
		new_records = records[0:index]
	default:
		new_records = records[0:index]
		new_records = append(new_records, records[(index + 1):]...)
	}

	return new_records
}



func printRecordsToStdout(records []Record) {
	printRecords(os.Stdout, records, "%v%v) %v\n%v%v\n")
}



func printRecordsToTempFile(records []Record, file *os.File) {
	printRecords(file, records, "%v%v) %v\n%vTags: %v\n\n")
}



func printRecords(out io.Writer, records []Record, format string) {
	// This is the number of records.
	m := len(records)
	// This is the number of digits in that number.
	n := len(strconv.FormatInt(int64(m), 10))

	spaces_bot := strings.Repeat(" ", (n + 2))

	for o := 0; o < m; o++ {
		spaces_top := ""

		if v := (n - len(strconv.FormatInt(int64(o + 1), 10))); v > 0 {
			spaces_top += strings.Repeat(" ", v)
		}

		fmt.Fprintf(out, format,
			spaces_top, (o + 1), records[o].Value,
			spaces_bot, strings.Join(records[o].Tags, ", "))
		// fmt.Fprintf(out
		//  "%v%v) %v%v\n%v%v\n",
		// 	spaces_top, (o + 1), records[o].MatchRate, records[o].Value,
		// 	spaces_bot, strings.Join(records[o].Tags, ", "))
	}
}



func promptForWantedRecord(verb string) string {
	fmt.Printf("%v%v these records: ", strings.ToUpper(string(verb[0])), string(verb[1:]))
	// fmt.Print("Enter the number(s) of the record(s) you want: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.TrimSpace(input)
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



func cleanInputTags(input string) []string {
	var clean []string
	ref := make(map[string]bool)

	parts := strings.Split(input, ",")

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)

		if _, in_ref := ref[trimmed]; in_ref == false {
			clean = append(clean, trimmed)
			ref[trimmed] = true
		}
	}

	return clean
}



func getWantedRecords(records []Record, input string) []Record {
	if strings.ToLower(input) == "all" {
		return records
	} else {
		ints := cleanWantedRecordsInput(input)

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



func makeRecordSelector(prompt_verb string, act func([]Record)) func([]Record) {
	selector := func(records []Record) {
		printRecordsToStdout(records)

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
	updater := func(records []Record) {
		updateWantedRecords(records)
		updateStoreFile(conf.Store, makeBackupUpdater(records, saveRecordToFile))
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
		updateStoreFile(conf.Store, makeBackupUpdater(records, dels))
	}

	return deleter
}



func makeBackupUpdater(records []Record, bkAct func(*os.File, Record)) func(*os.File) func(Record) {
	bk := func(bk_file *os.File) func(Record) {
		_bk := func(record Record) {
			if (len(records) == 0) {
				bk_file.WriteString(joinRecord(record))
			} else {
				bk_line := true

				for n, chk := range records {
					if ((chk.Value == record.Value) && (reflect.DeepEqual(chk.Tags, record.Tags))) {
						bkAct(bk_file, record)
						records = removeRecord(records, n)
						bk_line = false
						break
					}
				}

				if bk_line {
					bk_file.WriteString(joinRecord(record))
				}
			}
		}

		return _bk
	}

	return bk
}



func makeEditor(conf Config) func([]Record) {
	ed := func(records []Record) {
		// Make tmp file
		tmp_name := getTempFileName("edit")
		tmp_file := createFile(tmp_name)
		// defer tmp_file.Close()

		// Print instructions, records to tmp file
		tmp_file.WriteString(EditFileInstructions)
		printRecordsToTempFile(records, tmp_file)
		tmp_file.Close()

		// Open tmp file in editor
		ed := checkEditor(conf.Editor)
		pipeToToolAsArg(tmp_name, ed)

		// Wait for editor to close

		// Open temp file
		ed_recs := parseRecordsFromTempFile(tmp_name)
		mod_recs, new_recs := collateRecordsByIndex(records, ed_recs)

		fmt.Printf("Parsed records from temp file `%v`:\n%v\n%v\n", tmp_name, mod_recs, new_recs)

		bk := func(bk_file *os.File) func(Record) {
			_bk := func(record Record) {
				if (len(mod_recs) == 0) {
					bk_file.WriteString(joinRecord(record))
				} else {
					bk_line := true

					for n, mod := range mod_recs {
						if ((mod[0].Value == record.Value) && (reflect.DeepEqual(mod[0].Tags, record.Tags))) {
							bk_file.WriteString(joinRecord(mod[1]))

							var new_mods [][]Record
							switch {
							case n == 0:
								new_mods = mod_recs[1:]
							case n == (len(mod_recs) - 1):
								new_mods = mod_recs[0:n]
							default:
								new_mods = mod_recs[0:n]
								new_mods = append(new_mods, mod_recs[(n + 1):]...)
							}
							mod_recs = new_mods

							bk_line = false
							break
						}
					}

					if bk_line {
						bk_file.WriteString(joinRecord(record))
					}
				}
			}

			return _bk
		}

		updateStoreFile(conf.Store, bk)

		if len(new_recs) > 0 {
			appendRecordsToStore(conf.Store, new_recs)
		}

		err := os.Remove(tmp_name)
		checkForError(err)

		// updateStoreFile(conf.Store, mod_recs, saveRecordToFile)

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

	return ed
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



func parseRecordsFromTempFile(tmp_name string) map[int]Record {
	tmp_file, err := os.Open(tmp_name)
	checkForError(err)
	defer tmp_file.Close()

	reader := bufio.NewReader(tmp_file)

	records := make(map[int]Record)

	var index int
	var value string
	var tags []string
	pairing := false

	re_val := regexp.MustCompile("^[ ]*([0-9]+)\\)[ ]+(.+)$")
	re_tag := regexp.MustCompile("^[ ]*(?:Tags:[ ]*)(.+)$")

	for {
		line, last := readNextEntry(reader, '\n')

		if pairing && (line == "" || last) {
			record := Record{}
			record.Value = value
			record.Tags = tags
			records[index] = record
			pairing = false
		} else if n := re_val.FindStringSubmatch(line); n != nil {
			if pairing {
				record := Record{}
				record.Value = value
				record.Tags = tags
				records[index] = record
				pairing = false
			}

			chk, err := strconv.Atoi(n[1])
			checkForError(err)
			index = chk - 1
			value = strings.TrimSpace(n[2])
			pairing = true
		} else if n := re_tag.FindStringSubmatch(line); n != nil {
			tags = cleanInputTags(strings.TrimSpace(n[1]))
		}

		if last {
			break;
		}
	}

	return records
}



func collateRecordsByIndex(ref_recs []Record, new_recs map[int]Record) ([][]Record, []Record) {
	var collated [][]Record
	for index, old_rec := range ref_recs {
		new_rec, in := new_recs[index]

		if in {
			new_rec.Meta = old_rec.Meta
			collated = append(collated, []Record{old_rec, new_rec})
			delete(new_recs, index)
		}
	}

	var additions []Record
	if len(new_recs) > 0 {
		for _, record := range new_recs {
			record.Meta = []string{strconv.FormatInt(time.Now().Unix(), 10), "0", "0"}
			additions = append(additions, record)
		}
	}

	return collated, additions
}
