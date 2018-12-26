package main

import (
	"bytes"
	"fmt"
	"io"
	// "io/ioutil"
	"os"
	"os/exec"
	//"strings"
	"unicode/utf8"
)


// makeRecordPiper makes the Pipe search action function: the
// returned function will receive the slice of wanted Records and
// pipe the values of each to an external tool.
func makeRecordPiper(act string, caller func([]Record, string)) func([]Record) {
	piper := func(records []Record) {
		caller(records, act)
	}
	return piper
}

// pipeRecordsAsStdin pipes the given Records to the tool named by the
// given path on stdin.
func pipeRecordsAsStdin(records []Record, tool string) {
	cmd := exec.Command(tool)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic("Failed to open stdin pipe for records!")
	}

	defer func() {
		for _, record := range records {
			io.WriteString(stdin, record.Value + "\n")
		}
		stdin.Close() // Close the pipe, thereby sending EOF.
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
}

// pipeToToolAsArg pipes the given string to the tool named by the
// given path as an argument, so in the form `tool str`.
func pipeToToolAsArg(str string, tool string) {
	cmd := exec.Command(tool, str)
	printErr := func(err error) {
		fmt.Fprintf(os.Stderr, "Error running `%v %v`: %v\n", tool, str, err)
	}
	runCommand(cmd, printErr)
}

// runCommand runs the given Command and checks for an error. If an
// error occurs, it gets passed to the given handler function.
func runCommand(cmd *exec.Cmd, printErr func(error)) {
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err == nil {
		if utf8.RuneCountInString(out.String()) > 0 {
			fmt.Printf("OUTPUT? %v\n", out.String())
		}
	} else {
		printErr(err)
	}
}
