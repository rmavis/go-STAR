package main

import (
	"os"
)


// makeDeleter makes the Delete search action function: the returned
// function will receive the slice of wanted Records and ensure they
// are not included in the updated store file.
func makeDeleter(conf *Config) func([]Record) {
	dels := func(_ *os.File, _ Record) {
		// fmt.Printf("Not including entry in backup (%v)\n", record)
	}

	deleter := func(records []Record) {
		updateStoreFile(conf.Store, makeBackupUpdater(records, dels))
	}

	return deleter
}

// removeRecord returns a copy of the given slice of Records but
// without the element on the given index.
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
