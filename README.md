![Gorilla](https://i.imgur.com/lX7Vzr0.png)

# Gorilla Programming Language

[![Go Reference](https://pkg.go.dev/badge/github.com/SnowballSH/Gorilla.svg)](https://pkg.go.dev/github.com/SnowballSH/Gorilla)

---

## NEWS:

This is the unstable/developing branch of Gorilla 1.0. It is current work-in-progress.

Gorilla 1.0 is focusing on speed and safety. I am aiming for
100% [test coverage](https://app.codecov.io/gh/SnowballSH/Gorilla)

**If you are expecting a fully-working gorilla source code, visit
the [0.x branch](https://github.com/SnowballSH/Gorilla/tree/0.x)**

If you are a golang developer, and you know some basic runtime or parsing knowledge, you are more than welcome to
contribute!

If you are not, feel free to create issues about what you are expecting.

---

### Usage of the Gorilla library

Download:

```bash
go get github.com/SnowballSH/Gorilla
```

Basic usage:

```go
import "github.com/SnowballSH/Gorilla/exports"

code := "'Hello, world!'"

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