package main

import (
	"fmt"
	"strings"
)





func makeSearchAction(conf Config, action_code []int, terms []string) func() {
	// fmt.Printf("Action code: %v\n", action_code)

	match_act := getMatchAction(conf, action_code)
	match_lim := getMatchLim(conf, action_code, terms)
	matcher := makeMatcher(terms, match_lim)
	sorter := makeSorter(action_code)

	action := func() {
		records := readRecordsFromStore(conf.Store, matcher)
		sorter(records)
		match_act(records)
	}

	return action
}



func getMatchAction(conf Config, action_code []int) func([]Record) {
	var act func([]Record)

	switch {
	case action_code[0] == 4:  // Dump all.
		act = printRecordsToStdout

	case action_code[1] == 0:  // Read from config.
		switch {
		case conf.Action == "copy":  // "copy" and "open" are shortcuts.
			act = makeRecordSelector("copy", makeActAndUpdater(conf, pipeRecordsToPbcopy))
		case conf.Action == "open":
			act = makeRecordSelector("open", makeActAndUpdater(conf, pipeRecordsToOpen))
		case conf.Action == "browse" || conf.Action == "":
			act = printRecordsToStdout
		default:                     // Any external command can be specified.
			piper := makeRecordPiper(conf.Action)
			act = makeRecordSelector(conf.Action, makeActAndUpdater(conf, piper))
		}

	case action_code[1] == 1:  // pbcopy
		act = makeRecordSelector("copy", makeActAndUpdater(conf, pipeRecordsToPbcopy))

	case action_code[1] == 2:  // open
		act = makeRecordSelector("open", makeActAndUpdater(conf, pipeRecordsToOpen))

	case action_code[1] == 3:  // edit
		act = makeRecordSelector("edit", makeEditor(conf))

	case action_code[1] == 4:  // delete
		act = makeRecordSelector("delete", makeDeleter(conf))

	case action_code[1] == 5:  // browse
		act = printRecordsToStdout

	default:  // Bork.
		act = printRecordsToStdout
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
func makeMatcher(terms []string, lim int) func(Record) (float64, bool) {
	matcher := func(record Record) (float64, bool) {
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
		}

		if matches < lim {
			return 0.0, false;
		} else {
			match_rate := 0.0

			for o := 0; o < len(match_rates); o++ {
				match_rate += match_rates[o]
			}

			return (match_rate / float64(len(match_rates))), true;
			// return match_rate
		}
	}

	return matcher
}
