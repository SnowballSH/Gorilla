## Get started

### Installation

- Microsoft Windows 10
    - Go to the [latest release](https://github.com/SnowballSH/Gorilla/releases)
    - Scroll down, download `gorilla.exe`
    - If you trust me, ignore windows defender and keep it
        - this is not a malicious executable :)
    - Optional, add `gorilla.exe` to PATH
    
- Other OS
    - Due to the problem that I only have windows 10, 
    there is no pre-built binary available for other OS.
    - Download [Golang](https://golang.org/dl/) if you haven't
        - Test installation by `go version`
    - Download [GNU Makefile](https://www.gnu.org/software/make/) if you haven't
        - Test installation by `make -v`
        
    - Pull the latest commit
        
    - Get your binary by doing `make` in the root folder
    - Find `gorilla.app` or `gorilla` depending on your OS

Warning: All command-line code are for Microsoft Windows 10, Powershell.
It may or may not work in other terminals / OS.

### REPL

Simply do
```batch
gorilla
```
to fire up the gorilla REPL

type in random things to try!

### Run File

Make a custom file called `filename.gor`

Write some Gorilla code there

do
```batch
gorilla filename.gor
```
to run it

### Tutorial

View the [Tutorial](https://github.com/SnowballSH/Gorilla/blob/master/docs/tut.md)!
