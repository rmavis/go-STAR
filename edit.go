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
