## Welcome to the Gorilla Programming Language!

### What is Gorilla?

Gorilla is a tiny, dynamically typed, flexible programming language, written in [Golang](https://golang.org/)

Gorilla is built to make simple, flexible, well-understood programs in a relatively fast way.

### Why should I use Gorilla?

- It is **gorilla-crawling-blaze-ish fast** compared to other dynamic, interpreted languages
    - According to a fibonacci benchmark, Gorilla is **36% faster** than [Python](https://www.python.org/)
    - Gorilla is compiled into [Gorilla Bytecode](https://github.com/SnowballSH/Gorilla/blob/master/code/code.go)
- It is **dynamic and flexible**
    - You **don't have to worry** about all those `semicolons or modules` anymore!
    - Gorilla has **Python/Ruby style** assignment expression: `name = value`, which looks clean
    - Gorilla also has **concise operators** like `array <- item`, which obviously puts `item` to the end of `array`
- While it is dynamic, you can still write type-safe code
    - 99% of Gorilla's builtin functions are **type-safe**
    - Gorilla has **builtin support** for type annotations and assertions for functions

- ### [Learn Gorilla by Examples](https://snowballsh.github.io/Gorilla/byexample)
- ### [Gorilla Language Reference](https://snowballsh.github.io/Gorilla/api)

<br>

---

<br>

## Setting Up Gorilla

##### Warning: gorilla is still not production-ready, and it still contains some bugs, please report them if you found one!

<br>

### There are 3 ways to set up Gorilla:

1. If you are on **Windows 10** (and you trust I don't spread viruses), [download the latest binary executable](https://github.com/SnowballSH/Gorilla/releases)
2. Try **Online**: [Gorilla Playground](https://snowballsh.me/Gorilla-Playground/)
3. Download [Golang](https://golang.org/) (recommend v1.5), `git pull https://github.com/SnowballSH/Gorilla` and `go build`

<br>

## Hello, world!

The simplest way to write a Hello World program in Gorilla:

```ruby
println("Hello, world!")
# Hello, world!
```

#### Note: `print` prints without a new line and `println` prints with a new line

<br>

- ### [Learn Gorilla by Examples](https://snowballsh.github.io/Gorilla/byexample)
- ### [Gorilla Language Reference](https://snowballsh.github.io/Gorilla/api)
