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


//
// Functions for reading and acting on records in the store file.
//

// readRecordsFromFile reads the file named by the given string,
// parses each well-formed entry into a Record, and passes the Record
// to a function that determines whether it "matches". Each matching
// Record in added to a slice, and that slice of Records is retutned.
func readRecordsFromFile(file_name string, getMatchInfo func(Record) (float64, bool)) []Record {
	var records []Record

	act := func(record Record) {
		match_rate, matches := getMatchInfo(record)

		if matches {
			record.MatchRate = match_rate
			records = append(records, record)
		}
	}

	forEachRecordInFile(file_name, act)

	return records
}

// forEachRecordInFile reads the file named by the given string and,
// for each well-formed entry, it transforms the entry to a Record
// and passes the Record to the given function.
func forEachRecordInFile(file_name string, actOnRecord func(Record)) {
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

// updateStoreFile will "update" the file named by the given string
// by first making a backup and then renaming the backup over the
// original. The backup process is determined by the `bkMaker` param,
// which must be a function that returns a function that can be used
// on each record in the file.
func updateStoreFile(file_name string, bkMaker func(*os.File, Record)) {
	bk_name := file_name + "_bk_" + strconv.FormatInt(time.Now().Unix(), 10)
	bk_file, err := os.Create(bk_name)
	checkForError(err)
	defer bk_file.Close()

	updater := func(record Record) {
		bkMaker(bk_file, record)
	}
	forEachRecordInFile(file_name, updater)

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


//
// Utility functions.
//

// readNextEntry reads the given IO buffer up to the next separator.
// It returns the string read, removing whitespace and the separator.
func readNextEntry(reader *bufio.Reader, separator byte) (string, bool) {
	record, err := reader.ReadBytes(separator)
	last := false
	if err != nil {
		last = true
		// fmt.Printf("Error! %v (%v)\n", err, string(record))
	}
	return strings.TrimSpace(string(record)), last;
}

// saveRecordToFile writed the given Record to the given file.
func saveRecordToFile(file *os.File, record Record) {
	_, err := file.WriteString(joinRecord(record))
	checkForError(err)
}

// getTempFileName returns a temp file whose name includes the given
// string.
func getTempFileName(fx string) string {
	usr, err := user.Current()
	checkForError(err)
	return os.TempDir() + "/" + usr.Name + "_star_" + fx + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".tmp"
}

// createFile creates a file named by the given string.
func createFile(path string) *os.File {
	file, err := os.Create(path)
	checkForError(err)
	file.Chmod(0644)
	return file
}

// doesFileExist checks if a file named by the given string exists.
func doesFileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
