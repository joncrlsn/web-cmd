package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// SynchronizedFile synchronizes writing to the underlying file
type SynchronizedFile struct {
	file  *os.File
	mutex sync.Mutex
}

// NewSynchronizedFile synchronizes writing to a writer
func NewSynchronizedFile(f *os.File) *SynchronizedFile {
	sf := &SynchronizedFile{file: f}
	return sf
}

// WriteString writes to the file
func (sf *SynchronizedFile) WriteString(text string) (int, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	return sf.file.WriteString(text)
}

// Close closes the file
func (sf *SynchronizedFile) Close() error {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	return sf.file.Close()
}

// Commander is a struct that represents an exec.Cmd process which
// can take input (via the WriteString method) and processes all
// output (including the input) to a file.
type Commander struct {
	cmd     *exec.Cmd
	pIn     io.WriteCloser
	pOut    io.ReadCloser
	pErr    io.ReadCloser
	outFile *SynchronizedFile
	done    bool
	verbose bool
}

// NewCommander builds a Commander instance that will run a shell command and write it's output to a file
func NewCommander(command, outFile string, verbose bool) (*Commander, error) {

	c := &Commander{cmd: exec.Command(command), done: false, verbose: verbose}

	// Create the outWriter
	file, err := os.Create(outFile)
	if err != nil {
		return nil, err
	}
	c.outFile = NewSynchronizedFile(file)

	c.pIn, err = c.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	c.pOut, err = c.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.pErr, err = c.cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	//
	// Capture standard error and print it to outFile
	//
	go func() {
		defer c.pErr.Close()
		errReader := bufio.NewReader(c.pErr)
		errScanner := bufio.NewScanner(errReader)
		for errScanner.Scan() {
			line := errScanner.Text()
			if c.verbose {
				fmt.Println("stderr:", line)
			}
			c.outFile.WriteString(line + "\n")
			//fmt.Fprintln(c.outFile, line)
		}
		if err := errScanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading standard error from %s:", command, err)
		}
	}()

	//
	// Capture standard output and print it to a file
	//
	go func() {
		reader := bufio.NewReader(c.pOut)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Bytes()
			text := fmt.Sprintf("%s\n", massageOutputBytes(line))
			if c.verbose {
				fmt.Printf("Out: %s\n", text)
			}
			c.outFile.WriteString(text)
			//fmt.Fprintf(c.outFile, "%s\n", line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading standard output from %s:", command, err)
		}
	}()

	if c.verbose {
		fmt.Println("Created new Commander for command:", command)
	}
	return c, nil
}

// writeToStdIn can be overridden for custom behaviors.
// Like perhaps you do not want a newline written with the text
// for some strings.
var writeToStdIn = func(stdIn io.WriteCloser, text string) error {
	_, err := fmt.Fprintln(stdIn, text)
	return err
}

// massageOutputBytes can be overridden for custom behaviors.
// The bytes returned will be written to the output file.
var massageOutputBytes = func(bytes []byte) []byte {
	// This does nothing, but customizations may do something different
	return bytes
}

// WriteString writes to std in of the command and to the output file
func (c *Commander) WriteString(text string) (int, error) {
	if len(text) == 0 {
		return 0, nil
	}
	if c.verbose {
		fmt.Println("input:", text)
	}

	err := writeToStdIn(c.pIn, text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing to stdin:", err)
		return 0, err
	}

	return c.outFile.WriteString("> " + text + "\n")
}

// Close closes the output file
func (c *Commander) Close() error {
	c.done = true
	err := c.outFile.Close()
	if c.verbose {
		fmt.Println("Commander.Close() err:", err)
	}
	return err
}

// Start starts the exec.Cmd instance
func (c *Commander) Start() error {
	err := c.cmd.Start()
	if c.verbose {
		fmt.Println("Commander.Start() err:", err)
	}
	return err
}

// Wait waits for the exec.Cmd instance to finish and then closes
func (c *Commander) Wait() error {
	if c.verbose {
		fmt.Println("Commander.Wait() waiting...")
	}
	defer c.Close()
	err := c.cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error from Wait():", err)
	}
	return err
}
