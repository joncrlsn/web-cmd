package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

// main runs the bc command via a Commander instance
// and sends 4 calculations before quitting.
func Test_Commander(t *testing.T) {
	c, err := NewCommander("/usr/bin/bc", "tmp.txt", false /*true=verbose*/)
	if err != nil {
		panic(err)
	}
	if err := c.Start(); err != nil {
		panic(err)
	}

	// Calculate several things and then send "quit"
	go func() {
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("Calculating 4*5")
		if _, err := c.WriteString("4*5"); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("Calculating (2^3)^2")
		if _, err := c.WriteString("(2^3)^2"); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("Calculating 2^8")
		if _, err := c.WriteString("2^8"); err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 500)
		//fmt.Println("quit")
		if _, err := c.WriteString("quit"); err != nil {
			panic(err)
		}
	}()

	c.Wait()
	fmt.Println("Check for output in tmp.txt")

	// Read tmp.txt into line (which is a slice of strings)
	file, err := os.Open("tmp.txt")
	if err != nil {
		t.Fatal("Error opening tmp.txt", err)
	}
	defer file.Close()
	defer os.Remove("tmp.txt")

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Ensure the 6th line is correct
	if lines[5] != "256" {
		t.Fatal("Expected 256 in the 6th line of tmp.txt... instead got %s", lines[5])
	}
}
