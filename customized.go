package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// customize potentially redefines some functions used by commander
func customize(cmd string) {
	if strings.Contains(cmd, "pianobar") {
		writeToStdIn = writeToStdInPianobar
		massageOutputBytes = massageOutputBytesPianobar
	}
}

//
// Pianobar customizations
//

var writeToStdInPianobar = func(stdIn io.WriteCloser, text string) error {
	if text == "s" {
		fmt.Println("Writing to pianobar without newline: ", text)
		// User is selecting a channel, so write no newline
		_, err := fmt.Fprint(stdIn, text)
		return err
	}
	fmt.Println("Writing to pianobar with newline: ", text)
	_, err := fmt.Fprintln(stdIn, text)
	return err
}

var escLeftSquare2K = `\x1b\[2K` // x1b5b324b or <Esc>[2K
var pianobarRemoveRegex = regexp.MustCompile(escLeftSquare2K + ".*\r" + escLeftSquare2K)

var massageOutputBytesPianobar = func(bytesIn []byte) []byte {

	// first pass at clearning up the output
	bytes2 := pianobarRemoveRegex.ReplaceAll(bytesIn, []byte{})

	// 2nd pass at clearning up the output
	removePrefix := bytes.NewBufferString("[2K").Bytes()
	bytes3 := bytes.Replace(bytes2, removePrefix, []byte{}, 1)
	return bytes3
}
