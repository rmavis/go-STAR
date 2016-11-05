package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)





func makeCreateAction(conf Config, action_code []int, terms []string) func() {
	action := func() {
		record := makeRecordFromInput(terms)
		appendRecordToStore(conf.Store, record)
	}

	return action
}



func makeRecordFromInput(terms []string) Record {
	val := terms[0]
	tags := terms[1:]

	// t := time.Now()
	// u := t.Unix()
	s := strconv.FormatInt(time.Now().Unix(), 10)

	times := []string{s, "0", "0"}

	return Record{val, tags, times, 0.0}
}



func appendRecordToStore(file_name string, record Record) {
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY, 0600)
	checkForError(err)

	entry := joinRecord(record)

	_, err = file.WriteString(entry)
	checkForError(err)

	file.Close()
}



func joinRecord(record Record) string {
	parts := []string{
		record.Value,
		string(RecordSeparator),
		strings.Join(record.Tags, string(UnitSeparator)),
		string(RecordSeparator),
		strings.Join(record.Meta, string(UnitSeparator)),
		string(GroupSeparator),
		"\n",
	}

	return strings.Join(parts, "")
}
