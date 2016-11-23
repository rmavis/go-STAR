package main

import (
	"bufio"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)



// The value of EditFileInstructions will be written to the top of
// the temporary file used for editing entries.
const EditFileInstructions = `# STAR will read this file and update its store with the new values.
# 
# An entry in the store file includes two lines in this file, so STAR
# will expect to read lines from this file in pairs that look like:
#
#   1) http://settlement.arc.nasa.gov/70sArtHiRes/70sArt/art.html
#      Tags: art, NASA, space
#
# Those parts are:
# - At the start of a line (spaces excluded) a number followed by a closing parenthesis
# - The entry, being the string that gets copied, opened, etc
# - At the start of a line (spaces excluded) the word "Tags" followed by a colon
# - The tags, being a comma-separated list
#
# You can remove entries from the store file by deleting the line
# pairs, and you can add entries by creating more.
#
# Lines that start with a pound sign will be ignored.

`





// makeEditor returns the Edit search action function: the returned
// function will receive the slice of wanted Records and run the edit
// routine of printing the records to a temp file, reading & parsing
// that temp file, and incorporating the changes into the updated
// store file. Through this process records can be updated and added
// but not deleted.
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
		// fmt.Printf("Parsed records from temp file `%v`:\n%v\n%v\n", tmp_name, mod_recs, new_recs)

		updateStoreFile(conf.Store, makeEditUpdater(mod_recs))

		if len(new_recs) > 0 {
			appendRecordsToFile(conf.Store, new_recs)
		}

		err := os.Remove(tmp_name)
		checkForError(err)
	}

	return ed
}



// makeEditUpdater does the same thing as `makeBackupUpdater` except
// it operates on pairs of Records instead of individuals. For each
// pair, the first is the old entry, used for comparison, and the
// second is the new one, which gets saved in place of the old one.
func makeEditUpdater(rec_pairs [][]Record) func(*os.File) func(Record) {
	bk := func(bk_file *os.File) func(Record) {
		_bk := func(record Record) {
			if (len(rec_pairs) == 0) {
				saveRecordToFile(bk_file, record)
			} else {
				bk_line := true

				for n, mod := range rec_pairs {
					if ((mod[0].Value == record.Value) && (reflect.DeepEqual(mod[0].Tags, record.Tags))) {
						saveRecordToFile(bk_file, mod[1])

						var new_mods [][]Record
						switch {
						case n == 0:
							new_mods = rec_pairs[1:]
						case n == (len(rec_pairs) - 1):
							new_mods = rec_pairs[0:n]
						default:
							new_mods = rec_pairs[0:n]
							new_mods = append(new_mods, rec_pairs[(n + 1):]...)
						}
						rec_pairs = new_mods

						bk_line = false
						break
					}
				}

				if bk_line {
					saveRecordToFile(bk_file, record)
				}
			}
		}

		return _bk
	}

	return bk
}



// parseRecordsFromTempFile reads the file named by the given string
// and translates that data into a map of Records, which it returns.
// A map is used instead of a slice because, in the edit file, each
// record will be preceded by a number, as they are when printed in
// the terminal. That number corresponds to an index in the slice of
// wanted Records, and that's how the updates are paired with the
// existing Records.
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



// collateRecordsByIndex pairs the Records parsed from the edit file
// with the slice of wanted Records. If Records are present that do
// not correspond to the slice of wanted Records, then those are new
// Records, and they'll be added to the updated store file.
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



// cleanInputTags transforms the string of tags from the edit file
// into a slice of strings that can be used in the Record structure.
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
