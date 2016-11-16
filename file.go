package main

import (
	"bufio"
	// "fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)





func doesFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}



func createFile(path string) {
	file, err := os.Create(path)
	checkForError(err)

	file.Chmod(0644)
}



func updateStoreFile(file_name string, records []Record, bkAct func(*os.File, Record)) {
	bk_name := file_name + "_bk_" + strconv.FormatInt(time.Now().Unix(), 10)
	bk_file, err := os.Create(bk_name)
	checkForError(err)
	defer bk_file.Close()

	act := func(record Record) {
		if (len(records) == 0) {
			bk_file.WriteString(joinRecord(record))
		} else {
			x := len(records) - 1
			bk_line := true

			for n, chk := range records {
				if ((chk.Value == record.Value) && (reflect.DeepEqual(chk.Tags, record.Tags))) {
					bkAct(bk_file, chk)
					// bk_file.WriteString(joinRecord(chk))

					var new_records []Record
					switch {
					case n == 0:
						new_records = records[1:]
					case n == x:
						new_records = records[0:n]
					default:
						new_records = records[0:n]
						new_records = append(new_records, records[(n + 1):]...)
					}

					records = new_records
					bk_line = false
					break
				}
			}

			if bk_line {
				bk_file.WriteString(joinRecord(record))
			}
		}
	}

	forEachRecordInStore(file_name, act)

	// Get perms from store
	// f_info, err := os.Stat(file_name)
	// checkForError(err)

	// Set backup perms to store perms
	// err = bk_file.Chmod(f_info.Mode().Perm())
	// checkForError(err)

	// Rename backup
	err = os.Rename(bk_name, file_name)
	checkForError(err)
}



func forEachRecordInStore(file_name string, actOnRecord func(Record)) {
	file_handle, err := os.Open(file_name)
	checkForError(err)
	defer file_handle.Close()

	reader := bufio.NewReader(file_handle)

	for {
		entry, last := readNextEntry(reader)
		parts := splitEntry(entry)

		if (doesEntryHaveParts(parts)) {
			actOnRecord(makeRecordFromParts(parts))
		} else {
			// fmt.Printf("Record is missing components: %v\n", entry)
		}

		if last {
			break;
		}
	}
}
