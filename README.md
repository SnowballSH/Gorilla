![Gorilla](https://i.imgur.com/lX7Vzr0.png)

# Gorilla Programming Language

[![Go Reference](https://pkg.go.dev/badge/github.com/SnowballSH/Gorilla.svg)](https://pkg.go.dev/github.com/SnowballSH/Gorilla)

---

**If you are expecting a fully-working gorilla source code, visit
the [0.x branch](https://github.com/SnowballSH/Gorilla/tree/0.x)**

---

### About

Gorilla is a dynamic, interpreted programming language written in Go and Plan 9 Assembly.

It is made for creating fast, async, efficient, and simple apps.

---

### Usage of the Gorilla library

Download:

```bash
go get github.com/SnowballSH/Gorilla
```

Basic usage:

```go
import "github.com/SnowballSH/Gorilla/exports"

code := "'Hello,' + ' world!'"

// Compile Gorilla to bytecode
res, err := exports.CompileGorilla(code)

if err != nil {
	panic(err)
}

// Execute Gorilla bytecode along with the source text for debugging
vm, lastItem, err := exports.ExecuteGorillaBytecodeFromSource(res, code)

println(lastItem.Inspect()) // 'Hello, world!'
```

---

The following links are all based on gorilla 0.4.0.1:

#### View more on the [website](https://snowballsh.me/Gorilla/)

#### Download the [latest release](https://github.com/SnowballSH/Gorilla/releases)

#### Play it [Online](https://snowballsh.me/Gorilla-Playground/)