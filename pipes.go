package main

import (
	"bytes"
	"fmt"
	// "io"
	// "io/ioutil"
	// "os"
	"os/exec"
	"strings"
	"unicode/utf8"
)





// makeRecordPiper makes the Pipe search action function: the
// returned function will receive the slice of wanted Records and
// pipe the values of each to an external tool.
func makeRecordPiper(act string) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act)
	}

	return piper
}



// pipeToToolAsArg pipes the given string to the tool named by the
// given path as an argument, so in the form `tool str`.
func pipeToToolAsArg(str string, tool string) {
	cmd := exec.Command(tool, str)

	printErr := func(err error) {
		fmt.Printf("Error running `%v %v`: %v\n", tool, str, err)
	}

	runCommand(cmd, printErr)
}



// pipeToToolAsStdin pipes the given string to the tool named by the
// given path as stdin, so in the form `str | tool`.
func pipeToToolAsStdin(str string, tool string) {
	cmd := exec.Command(tool)
	cmd.Stdin = strings.NewReader(str)

	printErr := func(err error) {
		fmt.Printf("Error running `%v | %v`: %v\n", str, tool, err)
	}

	runCommand(cmd, printErr)
}



// pipeRecordsToPbcopy is a convenience function for passing each of
// the given Records to `pbcopy` as stdin, which `pbcopy` requries.
func pipeRecordsToPbcopy(records []Record) {
	// pipeRecordsToExternalTool(records, PbcopyPath);

	for _, r := range records {
		pipeToToolAsStdin(r.Value, PbcopyPath)
	}
}



// pipeRecordsToOpen is a convenience function for passing each of
// the given Records to `open` as an argument, as `open` requires.
func pipeRecordsToOpen(records []Record) {
	pipeRecordsToExternalTool(records, OpenPath);
}



// pipeRecordsToExternalTool pipes each of the given Records as an
// argument to the tool named by the given path.
func pipeRecordsToExternalTool(records []Record, tool string) {
	for _, r := range records {
		pipeToToolAsArg(r.Value, tool)
	}
}



// runCommand runs the given Command and checks for an error. If an
// error occurs, it gets passed to the given handler function.
func runCommand(cmd *exec.Cmd, printErr func(error)) {
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err == nil {
		if utf8.RuneCountInString(out.String()) > 0 {
			fmt.Printf("%v\n", out.String())
		}
	} else {
		printErr(err)
	}
}



// func captureStdout(f func()) string {
// 	_stdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w

// 	f()

// 	w.Close()
// 	os.Stdout = _stdout

// 	var buf bytes.Buffer
// 	io.Copy(&buf, r)

// 	return buf.String()
// }
