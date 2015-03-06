//
// Copyright (c) 2015 Jon Carlson.  All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
//
// web-cmd webifies a basic interactive shell command (like bc or pianobar)
package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var version = "0.4"
var context = "/"
var port string
var command string
var verbose bool
var commander *Commander
var tempFileName string
var bodyFormat = ` <!doctype html> <html> <body> <h1>%s</h1> <a href="#bottom">bottom</a> <pre><code>%s</code></pre>
    <a name="bottom">
    <form method="post">
        <input name="cmd" placeholder="Enter a command" autofocus>
        <input type="submit" value="Submit">
    </form>
</body>
</html>`

func main() {

	var verFlag bool
	var helpFlag bool
	//var portStr string
	var err error
	flag.StringVarP(&command, "command", "c", "", "shell command to run (try bc on OSX and Linux)")
	flag.StringVarP(&port, "port", "p", "8080", "web port")
	flag.BoolVarP(&verFlag, "version", "V", false, "displays version information")
	flag.BoolVarP(&verbose, "verbose", "v", false, "prints extra information")
	flag.BoolVarP(&helpFlag, "help", "?", false, "displays usage help")
	flag.Parse()

	if verFlag {
		fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], version)
		fmt.Fprintln(os.Stderr, "Copyright (c) 2015 Jon Carlson.  All rights reserved.")
		fmt.Fprintln(os.Stderr, "Use of this source code is governed by the MIT license")
		fmt.Fprintln(os.Stderr, "that can be found here: http://opensource.org/licenses/MIT")
		os.Exit(1)
	}

	if helpFlag {
		usage()
	}

	if len(command) == 0 {
		fmt.Fprintln(os.Stderr, "The command flag (-c or --command) is required")
		os.Exit(1)
	}

	if verbose {
		fmt.Println("in verbose mode")
	}

	// Create a temp file name
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	tempFileName = filepath.Join(os.TempDir(), "web-cmd_"+hex.EncodeToString(randBytes))

	customize(command)
	if verbose {
		fmt.Println("Temporary output file:", tempFileName)
	}

	commander, err = startCommander(command, tempFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	portStr := ":" + port
	fmt.Printf("Listening on port %s\n", portStr)
	http.HandleFunc(context, handler)
	err = http.ListenAndServe(portStr, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting web server: %s", err)
		commander.Close()
	}
}

// handler returns the latest output
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGet(w, r)
	} else if r.Method == "POST" {
		handlePost(w, r)
	}
}

// handlGet returns the latest output from Commander
func handleGet(w http.ResponseWriter, r *http.Request) {
	if commander.done {
		errStr := "Commander is done.  Type 'restart' to resurrect it"
		fmt.Fprintf(w, bodyFormat, strings.Title(command), errStr)
		return
	}

	body, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		errStr := fmt.Sprintf("Error occurred reading file %s: %s", tempFileName, err)
		fmt.Fprintf(w, bodyFormat, strings.Title(command), errStr)
		return
	}
	fmt.Fprintf(w, bodyFormat, strings.Title(command), body)
}

// handlePost handles the stdin submission
func handlePost(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("cmd")
	if cmd == "restart" {
		if commander.done {
			fmt.Printf("Restarting %s commander\n", command)
			// create a new one
			c, err := startCommander(command, tempFileName)
			commander = c
			if err != nil {
				fmt.Printf("Error restarting %s: %s\n", command, err)
			}
		}
		http.Redirect(w, r, context, 302)
		return
	}
	commander.WriteString(cmd)
	time.Sleep(time.Second)
	http.Redirect(w, r, context, 302)
}

// startCommander populates the commander and tempFileName variables
func startCommander(shellCommand, tmpFileName string) (*Commander, error) {

	// Create a Commander that writes to the temp file
	c, err := NewCommander(shellCommand, tmpFileName, verbose)
	if err != nil {
		return nil, err
	}

	if err := c.Start(); err != nil {
		return nil, err
	}

	//go func() {
	//defer os.Remove(tmpFileName)
	//c.Wait()
	//}()

	return c, nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-?|--help] [-v|--version] -c <command>\n", os.Args[0])
	fmt.Fprintln(os.Stderr, `
Program flags are:
  -V, --version      : prints the version of web-cmd being run
  -?, --help         : prints a summary of the options accepted by web-cmd
  -v, --verbose      : prints more information about what is happening behind the scenes
  -c, --command      : shell command to be run
  -p, --port         : port web interface should listen on
`)

	os.Exit(2)
}
