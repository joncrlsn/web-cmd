# cmdweb - webifies an interactive OS command

The "itch" cmdweb scratches is this: Remotely control the pianobar app on a Raspberry Pi.  But it has other potential command-line webifying uses without exposing your entire shell over the web.  The other advantage of cmdweb over a web-based shell is that wherever this is accessed from, the same pianobar application instance is displayed and accessible. 

### download 
[osx](https://github.com/joncrlsn/pgrun/raw/master/bin-osx32/pgrun "OSX version")
[linux](https://github.com/joncrlsn/pgrun/raw/master/bin-linux32/pgrun "Linux version")
[windows](https://github.com/joncrlsn/pgrun/raw/master/bin-win32/pgrun.exe "Windows version")

### example
	cmdweb -c pianobar -p 8888  (then open your browser to http://localhost:8888)

#### flags/options (these mostly match psql arguments):
program flag/option  | explanation
-------------------: | -------------
  -V, --version      | prints the version of cmdweb being run
  -v, --verbose      | prints extra information to standard out
  -?, --help         | prints a summary of the program options
  -c, --command      | specifies the command you want to run
  -p, --port         | port webserver will listen on (default is 8080)

### todo
1. Need help/patches from others here
