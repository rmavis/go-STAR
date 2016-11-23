package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)





func makeActAndUpdater(conf Config, act func([]Record)) func([]Record) {
	updater := func(records []Record) {
		updateWantedRecords(records)
		updateStoreFile(conf.Store, makeBackupUpdater(records, saveRecordToFile))
		act(records)
	}

	return updater
}



// updateWantedRecords updates the metadata for each Record in the
// given slice.
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
