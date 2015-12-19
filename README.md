# web-cmd - webifies an interactive OS command

The "itch" web-cmd scratches is this: Remotely control the pianobar app on a Raspberry Pi.  But it has other potential command-line webifying uses without exposing your entire shell over the web.  The other advantage of web-cmd over a web-based shell is that wherever this is accessed from, the same pianobar application instance is displayed and accessible. 

### example
	web-cmd -c pianobar -p 8888  (then open your browser to http://localhost:8888)

#### options
        option  | explanation
--------------: | -------------
  -V, --version | prints the version of web-cmd being run
  -v, --verbose | prints extra information to standard out
  -?, --help    | prints a summary of the program options
  -c, --command | specifies the command you want to run
  -p, --port    | port webserver will listen on (default is 8080)

#### todo
1. Add binaries?
2. Add help command that explains things like 'quit', 'restart', 'start', etc
3. When pianobar is paused too long, a stream error occurs and the file is never updated: "Error reading standard output from"

