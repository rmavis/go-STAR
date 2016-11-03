package main

import (
	"os"
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
