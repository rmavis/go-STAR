package main

import (
	"bufio"
	// "fmt"
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
func makeEditor(conf *Config) func([]Record) {
	ed := func(records []Record) {
		// Create the temp file, add the instructions and records.
		tmp_name := getTempFileName("edit")
		tmp_file := createFile(tmp_name)
		tmp_file.WriteString(EditFileInstructions)
		listRecordsToTempFile(records, tmp_file)
		tmp_file.Close()

		// Open temp file in the user's editor, wait for editor to close.
		ed := checkEditor(conf.Editor, getEnv("EDITOR", DefaultEditorPath))
		pipeToToolAsArg(tmp_name, ed)

		// Read and parse temp file.
		ed_recs := parseRecordsFromTempFile(tmp_name)
		edits, adds, dels := collateRecordsByIndex(records, ed_recs)
		// fmt.Printf("Parsed records from temp file `%v`:\nEDITS: %v\nNEWS: %v\nDELETIONS: %v\n", tmp_name, edits, adds, dels)

		// Delete the temp file.
		err := os.Remove(tmp_name)
		checkForError(err)

		// Update the store file with all those changes.
		saveEditsToStore(conf, adds, edits, dels)
	}

	return ed
}

// saveEditsToStore receives the user's config and three slices that
// contain records to add, edit, and delete, and does those things
// to the store file indicated in the config.
func saveEditsToStore(conf *Config, adds []Record, edits [][]Record, dels []Record) {
	editer := func(bk_file *os.File, record Record) {
		should_bk := true

		for n, mod := range edits {
			if ((mod[0].Value == record.Value) && (reflect.DeepEqual(mod[0].Tags, record.Tags))) {
				saveRecordToFile(bk_file, mod[1])
				edits = removeRecordPair(edits, n)
				should_bk = false
				break
			}
		}

		for n, del := range dels {
			if ((del.Value == record.Value) && (reflect.DeepEqual(del.Tags, record.Tags))) {
				dels = removeRecord(dels, n)
				should_bk = false
				break
			}
		}

		if should_bk {
			saveRecordToFile(bk_file, record)
		}
	}

	updateStoreFile(conf.Store, editer)
	if len(adds) > 0 {
		appendRecordsToFile(conf.Store, adds)
	}
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
func collateRecordsByIndex(ref_recs []Record, new_recs map[int]Record) ([][]Record, []Record, []Record) {
	var deletions []Record
	var collated [][]Record
	for index, old_rec := range ref_recs {
		new_rec, in := new_recs[index]
		if in {
			if ((new_rec.Value != old_rec.Value) || (!reflect.DeepEqual(new_rec.Tags, old_rec.Tags))) {
				new_rec.Meta = old_rec.Meta
				collated = append(collated, []Record{old_rec, new_rec})
			}
			delete(new_recs, index)
		} else {
			deletions = append(deletions, old_rec)
		}
	}

	var additions []Record
	if len(new_recs) > 0 {
		for _, record := range new_recs {
			record.Meta = []string{strconv.FormatInt(time.Now().Unix(), 10), "0", "0"}
			additions = append(additions, record)
		}
	}

	return collated, additions, deletions
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
