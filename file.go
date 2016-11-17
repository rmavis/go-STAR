package main

import (
	"bufio"
	// "fmt"
	"os"
	"os/user"
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



func updateStoreFile(file_name string, bkMaker func(*os.File) func(Record)) {
	bk_name := file_name + "_bk_" + strconv.FormatInt(time.Now().Unix(), 10)
	bk_file, err := os.Create(bk_name)
	checkForError(err)
	defer bk_file.Close()

	updater := bkMaker(bk_file)

	forEachRecordInStore(file_name, updater)

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



func saveRecordToFile(file *os.File, record Record) {
	_, err := file.WriteString(joinRecord(record))
	checkForError(err)
}
