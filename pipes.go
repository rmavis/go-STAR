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





func makeRecordPiper(act string) func([]Record) {
	piper := func(records []Record) {
		pipeRecordsToExternalTool(records, act)
	}

	return piper
}



func pipeToToolAsArg(str string, tool string) {
	cmd := exec.Command(tool, str)

	printErr := func(err error) {
		fmt.Printf("Error running `%v %v`: %v\n", tool, str, err)
	}

	runCommand(cmd, printErr)
}



func pipeToToolAsStdin(str string, tool string) {
	cmd := exec.Command(tool)
	cmd.Stdin = strings.NewReader(str)

	printErr := func(err error) {
		fmt.Printf("Error running `%v | %v`: %v\n", str, tool, err)
	}

	runCommand(cmd, printErr)
}



func pipeRecordsToPbcopy(records []Record) {
	// pipeRecordsToExternalTool(records, PbcopyPath);

	for _, r := range records {
		pipeToToolAsStdin(r.Value, PbcopyPath)
	}
}



func pipeRecordsToOpen(records []Record) {
	pipeRecordsToExternalTool(records, OpenPath);
}



func pipeRecordsToExternalTool(records []Record, tool string) {
	for _, r := range records {
		pipeToToolAsArg(r.Value, tool)
	}
}



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
