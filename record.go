package main

import (
	"strings"
)



// Record is a structure that contains each part of an entry.
// An entry is built from a Record.
// A Record is built from a well-formed entry.
// A well-formed entry has the structure specified in `joinRecord`.
type Record struct {
	Value string
	Tags []string
	Meta []string
	MatchRate float64
}



// These are the ASCII Group, Record, and Unit separator characters.
// They are of type `rune`.
const GroupSeparator = ''   // Separates records
const RecordSeparator = ''  // Separates parts of records (the value from the tags)
const UnitSeparator = ''    // Separates parts of parts (each tag from each tag)





// joinRecord receives a Record and returns a string. The Record's
// tags and metadata will be joined with with unit separator; the
// resulting strings will be joined to the value with the record
// separator; and the group separator and a newline will be appended
// to the resulting string. The return will be a well-formed entry.
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



// splitEntry receives a string and returns a slice of strings. The
// `entry` should still contain the trailing group separator (which
// splits entries). The string will be split on the record separator.
// If the entry is well-formed, the return will have three parts.
func splitEntry(entry string) []string {
	fields := strings.Split(strings.TrimSuffix(entry, string(GroupSeparator)), string(RecordSeparator))
	return fields
}



// splitField receives a string and returns a slice of strings. It
// splits on the unit separator. This is useful for creating, e.g.,
// a slice of tags.
func splitField(field string) []string {
	parts := strings.Split(field, string(UnitSeparator))
	return parts
}



// makeRecordFromParts receives a slice of strings and returns a
// Record. The slice should be a well-formed entry: a string, and two
// a lists of strings joined by the unit separator.
func makeRecordFromParts(entry []string) Record {
	return Record{entry[0], splitField(entry[1]), splitField(entry[2]), 0.0}
}



// doesEntryHaveParts receives a slice of strings and returns a bool
// indicating whether the slice contains three parts. A well-formed
// entry has three parts: the value, tags, and metadata.
func doesEntryHaveParts(entry []string) bool {
	if (len(entry) == 3) {
		return true
	} else {
		return false
	}
}
