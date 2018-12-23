package main

import (
	"fmt"
	"os"
	"strings"
)


// makeSearchAction returns a multi-part action function that follow
// the pattern of: read records from file, search those records with
// the given terms, print the matches, prompt for the wanted records,
// then act on those wanted records. The final action taken on the
// wanted records is determined by the given action code.
func makeSearchAction(conf Config, action_code []int, terms []string) func() {
	act := mergeConfigActions(&conf, action_code)
	// fmt.Printf("Final Action code: %v\n", action)

	match_act := getMatchAction(&conf, act)
	match_lim := getMatchLim(act, len(terms))
	matcher := makeMatcher(terms, match_lim)
	sorter := makeSorter(act, (len(terms) > 0))

	action := func() {
		records := readRecordsFromFile(conf.Store, matcher)
		sorter(records)
		match_act(records)
	}

	return action
}


// getMatchAction returns a function that acts on a slice of Records.
// This function will be the the final action taken on the wanted
// records as specified in the multi-part search action function. The
// user's config and the action code are required to create the
// context/scope for the final action.
func getMatchAction(conf *Config, action_code []int) func([]Record) {
	var act func([]Record)

	var printer func([]Record)
	if (action_code[4] == 1) {
		printer = dumpRecordValuesToStdout
	} else {
		printer = listRecordsToStdout
	}

	switch {
	case action_code[1] == 1:  // View, no select.
		act = makeRecordPrintCaller(printer)
	case action_code[1] == 2:  // Select and pipe.
		piper := makeRecordPiper(conf.Action, pipeToToolAsStdin)
		act = makeRecordSelector("pipe", printer, makeActAndUpdater(conf, piper))
	case action_code[1] == 3:  // edit
		act = makeRecordSelector("edit", printer, makeEditor(conf))
	case action_code[1] == 4:  // delete
		act = makeRecordSelector("delete", printer, makeDeleter(conf))
	default:  // Bork.
		fmt.Fprintf(os.Stderr, "Unrecognized action `%v`", action_code[1])
		act = makeRecordPrintCaller(printer)
	}

	return act
}


// getMatchLim returns an integer that specifies the number of
// matches that must occur between the given terms and the scanned
// Records for a Record to "match" the terms. This value can depend
// on the action code or the number of terms.
func getMatchLim(action_code []int, num_terms int) int {
	var lim int

	switch {
	case num_terms == 0:
		lim = 0
	case action_code[2] == 1:  // loose
		lim = 1
	case action_code[2] == 2:  // strict
		lim = num_terms
	default:  // Bork.
		fmt.Fprintf(os.Stderr, "Invalid match code (%v). Using '1'.\n", action_code[2])
		lim = 1
	}

	return lim
}


// makeMatcher returns a function that can be called in the record-
// reading process to determine if the read record "matches".
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
