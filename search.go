package main

import (
	"fmt"
	"strings"
)





func makeSearchAction(conf Config, action_code []int, terms []string) func() {
	match_act := getMatchAction(conf, action_code)
	match_lim := getMatchLim(conf, action_code, terms)
	matcher := makeMatcher(terms, match_lim)
	sorter := makeSorter(conf, action_code, terms)

	action := func() {
		records := readRecords(matcher, sorter)
		match_act(records)
	}

	return action
}



func getMatchAction(conf Config, action_code []int) func([]Record) {
	var act func([]Record)

	switch {
	case action_code[0] == 4:  // Dump all.
		act = printRecords

	case action_code[1] == 0:  // Read from config.
		switch {
		case conf.Action == "copy":  // "copy" and "open" are shortcuts.
			act = makeRecordSelector(pipeRecordsToPbcopy)
		case conf.Action == "open":
			act = makeRecordSelector(pipeRecordsToOpen)
		case conf.Action == "browse" || conf.Action == "":
			act = printRecords
		default:                     // Any external command can be specified.
			piper := makeRecordPiper(conf.Action)
			act = makeRecordSelector(piper)
		}

	case action_code[1] == 1:  // pbcopy
		act = makeRecordSelector(pipeRecordsToPbcopy)

	case action_code[1] == 2:  // open
		act = makeRecordSelector(pipeRecordsToOpen)

	case action_code[1] == 3:  // edit
		act = makeRecordSelector(editRecords)

	case action_code[1] == 4:  // delete
		act = makeRecordSelector(deleteRecords)

	case action_code[1] == 5:  // browse
		act = printRecords

	default:  // Bork.
		act = printRecords
	}

	return act
}



func getMatchLim(conf Config, action_code []int, terms []string) int {
	var lim int

	switch {
	case action_code[0] == 4:  // All.
		lim = 0
	case action_code[2] == 0:  // Read from config.
		if conf.FilterMode == "loose" {
			lim = 1
		} else {
			lim = len(terms)
		}
	case action_code[2] == 1:  // loose
		lim = 1
	case action_code[2] == 2:  // strict
		lim = len(terms)
	default:  // Bork.
		fmt.Printf("Invalid match code (%v). Using '1'.\n", action_code[2])
		lim = 1
	}

	return lim
}



// The idea for the matcher function is to check the value and the
// tags for each argument. For each match, a match rate will be
// added to a collection. If the number of matches is greater than
// the given limit, then the aggregate of the rates will be returned
// as the record's overall match rate. The aggregate is used instead
// of the average (sum of match rates / number of matches) because
// it makes sense that a record that matches multiple times should
// rate higher than those that don't.
func makeMatcher(terms []string, lim int) func(Record) float64 {
	matcher := func(record Record) float64 {
		var match_rates []float64
		matches := 0

		strs := make([]string, 1, len(record.Tags) + 1)
		strs[0] = record.Value
		for o := 0; o < len(record.Tags); o++ {
			strs = append(strs, record.Tags[o])
		}
		str_agg := strings.Join(strs, string(UnitSeparator))
		// fmt.Printf("Aggregate line: %v\n", str_agg)

		for o := 0; o < len(terms); o++ {
			mult := strings.Count(str_agg, terms[o])

			if mult == 0 {
				match_rates = append(match_rates, 0.0)
			} else {
				matches += 1
				match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) * float64(mult)) / float64(len([]rune(str_agg)))))
			}

			// fmt.Printf("Checking %v\n", record)

			// if strings.Contains(record.Value, terms[o]) {
			// 	// fmt.Printf("V %v contains %v\n", record.Value, terms[o])
			// 	matches += 1
			// 	match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) / float64(len([]rune(record.Value)))) * float64(strings.Count(record.Value, terms[o]))))
			// } else {
			// 	match_rates = append(match_rates, 0.0)
			// }

			// for i := 0; i < len(record.Tags); i++ {
			// 	if strings.Contains(record.Tags[i], terms[o]) {
			// 		matches += 1
			// 		// fmt.Printf("T %v contains %v\n", record.Tags[i], terms[o])
			// 		match_rates = append(match_rates, ((float64(len([]rune(terms[o]))) / float64(len([]rune(record.Tags[i])))) * float64(strings.Count(record.Tags[i], terms[o]))))
			// 	} else {
			// 		match_rates = append(match_rates, 0.0)
			// 	}
			// }
		}

		if matches < lim {
			return 0.0
		} else {
			match_rate := 0.0

			for o := 0; o < len(match_rates); o++ {
				match_rate += match_rates[o]
			}

			return (match_rate / float64(len(match_rates)))
			// return match_rate
		}
	}

	return matcher
}
