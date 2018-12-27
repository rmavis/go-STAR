package main

import (
	"os"
	"reflect"
)


// makeDeleter makes the Delete search action function: the returned
// function will receive the slice of wanted Records and ensure they
// are not included in the updated store file.
func makeDeleter(conf *Config) func([]Record) {
	deleter := func(records []Record) {
		saveDeletionsToStore(conf, records)
	}

	return deleter
}

func saveDeletionsToStore(conf *Config, records []Record) {
	deleter := func(bk_file *os.File, record Record) {
		should_bk := true

		for n, chk := range records {
			if ((chk.Value == record.Value) && (reflect.DeepEqual(chk.Tags, record.Tags))) {
				records = removeRecord(records, n)
				should_bk = false
				break
			}
		}

		if should_bk {
			saveRecordToFile(bk_file, record)
		}
	}

	updateStoreFile(conf.Store, deleter)
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

func removeRecordPair(pairs [][]Record, index int) [][]Record {
	var new_records [][]Record

	switch {
	case index == 0:
		new_records = pairs[1:]
	case index == (len(pairs) - 1):
		new_records = pairs[0:index]
	default:
		new_records = pairs[0:index]
		new_records = append(new_records, pairs[(index + 1):]...)
	}

	return new_records
}
