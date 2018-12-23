package main

import (
	"bytes"
	"fmt"
	"io"
	// "io/ioutil"
	// "os"
	"os/exec"
	//"strings"
	"unicode/utf8"
)


// makeRecordPiper makes the Pipe search action function: the
// returned function will receive the slice of wanted Records and
// pipe the values of each to an external tool.
func makeRecordPiper(act string, caller func(string, string)) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act, caller)
	}
	return piper
}

// pipeRecordsToExternalTool pipes each of the given Records as an
// argument to the tool named by the given path.
func pipeRecordsToExternalTool(records []Record, tool string, caller func(string, string)) {
	for _, r := range records {
		caller(r.Value, tool)
	}
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
	fmt.Printf("WOULD PIPE `%v` AS STDIN TO `%v`", str, tool)

	// cmd := exec.Command(tool)
	// cmd.Stdin = strings.NewReader(str)
	// cmd_err := func(err error) {
	// 	fmt.Printf("Error running `%v | %v`: %v\n", str, tool, err)
	// }
	// runCommand(cmd, cmd_err)


    cmd := exec.Command(tool)
    stdin, err := cmd.StdinPipe()
    if err != nil {
		fmt.Printf("Error (1) : %v\n", err)
    }
	defer stdin.Close()
    _, err = io.WriteString(stdin, str + "\n")
    if err != nil {
		fmt.Printf("Error (3) : %v\n", err)
    }
	cmd.Run()

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
