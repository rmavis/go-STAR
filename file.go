package main

import (
	"bufio"
	// "fmt"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"time"
)





func getTempFileName(fx string) string {
	usr, err := user.Current()
	checkForError(err)

	return os.TempDir() + "/" + usr.Name + "_star_" + fx + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".tmp"
}



func createFile(path string) *os.File {
	file, err := os.Create(path)
	checkForError(err)

	file.Chmod(0644)

	return file
}



func doesFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}



func readNextEntry(reader *bufio.Reader, separator byte) (string, bool) {
	record, err := reader.ReadBytes(separator)
	last := false

	if err != nil {
		last = true
		// fmt.Printf("Error! %v (%v)\n", err, string(record))
	}

	return strings.TrimSpace(string(record)), last;
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
			bk_line := true

			for n, chk := range records {
				if ((chk.Value == record.Value) && (reflect.DeepEqual(chk.Tags, record.Tags))) {
					bkAct(bk_file, chk)
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

	forEachRecordInStore(file_name, act)

	// If there are still records in `records`, will want to append those.  #TODO

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
		entry, last := readNextEntry(reader, GroupSeparator)
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
