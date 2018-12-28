package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)


// makeActAndUpdater returns a procedure for use in the Search action
// function in which the wanted Records will be updated and those
// updates will be written to the store file.
func makeActAndUpdater(conf *Config, act func([]Record)) func([]Record) {
	updater := func(records []Record) {
		updateRecordsMetadata(records)
		saveUpdatesToStore(conf, records)
		act(records)
	}

	return updater
}

// updateRecordsMetadata updates the metadata for each Record in the
// given slice.
func updateRecordsMetadata(records []Record) {
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

// saveUpdatesToStore writes the updated records to the user's store
// file.
func saveUpdatesToStore(conf *Config, records []Record) {
	updater := func(bk_file *os.File, record Record) {
		should_bk := true

		for n, chk := range records {
			if ((chk.Value == record.Value) && (reflect.DeepEqual(chk.Tags, record.Tags))) {
				saveRecordToFile(bk_file, chk)
				records = removeRecord(records, n)
				should_bk = false
				break
			}
		}

		if should_bk {
			saveRecordToFile(bk_file, record)
		}
	}

	updateStoreFile(conf.Store, updater)
}
