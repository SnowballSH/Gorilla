## Welcome to the Gorilla Programming Language!

### What is Gorilla?

Gorilla is a tiny, dynamically typed, flexible programming language, written in [Golang](https://golang.org/)

Gorilla is built to make simple, flexible, well-understood programs in a relatively fast way.

### Why should I use Gorilla?

- It is **gorilla-crawling-blaze-ish fast** compared to other dynamic, interpreted languages
    - According to a fibonacci benchmark, Gorilla is **36% faster** than [Python](https://www.python.org/)
    - Gorilla is compiled into [Gorilla Bytecode](https://github.com/SnowballSH/Gorilla/blob/master/code/code.go#L6-L100)
- It is **dynamic and flexible**
    - You **don't have to worry** about all those `semicolons or modules` anymore!
    - Gorilla has **Python/Ruby style** assignment expression: `name = value`, which looks clean
    - Gorilla also has **concise operators** like `array <- item`, which obviously puts `item` to the end of `array`
- While it is dynamic, you can still write type-safe code
    - 99% of Gorilla's builtin functions are **type-safe**
    - Gorilla has **builtin support** for type annotations and assertions for functions

<br>

## Setting Up Gorilla

##### Warning: gorilla is still not production-ready, and it still contains some bugs, please report them if you found one!

<br>

### There are 3 ways to set up Gorilla:

1. If you are on **Windows 10** (and you trust I don't spread viruses), [download the latest binary executable](https://github.com/SnowballSH/Gorilla/releases)
2. Try **Online**: [Gorilla Playground](https://snowballsh.me/Gorilla-Playground/) (made using Web Assembly)
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

## Learn Gorilla by Examples

### Greeting

Gorilla Version: 0.4.0-alpha+

We will create a function that greets everyone!

```go
func greet(name) {
    println("Hello, " + name + "!")
}
```

The `greet` function accepts a parameter `name` and prints `Hello, #{name}!`

Because Gorilla has **flexible syntax**, `{ ... }` is optional if the function **only contains 1 statement**

The above code is equal to:

```go
func greet(name)
    println("Hello, " + name + "!")
```

We can call the function by doing:

```ruby
greet("world")
# Hello, world!

greet("Sam")
# Hello, Sam!
```

Full code:

```go
func greet(name)
    println("Hello, " + name + "!")

greet("world")
greet("Sam")
```

Feel free to modify and try other things out!

<br>

### Number Guessing Game

Gorilla Version: 0.4.0+

Must be on local system!

[Full Code](https://github.com/SnowballSH/Gorilla/blob/master/examples/guess.gor)

Let's start with generating a random integer.

```rust
use "random"

number = random.intRange(100) + 1
```

The `use` keyword imports a module, in this case, `random`

`random.intRange(x)` generates a random number from 0 to x-1, so we need to add 1 to make it 1-based.

A number guessing name must have an infinite game loop, so:

```ruby
run = true
while run {
    guess = input("Guess a number between 1 and 100: ")
}
```

Here we let the user guess a number friendlily.

`input()` gets a line of user input!

To check if a string is an integer, use the `string.isInt()` function.

```ruby
while run {
    guess = input("Guess a number between 1 and 100: ")

    if !guess.isInt() {
        println("Your input is not an integer!")
        next
    }

    guess = guess.toInt()

    if guess > number
        println("Guess smaller!")
    else if guess < number
        println("Guess larger!")
    else {
        println("You guess the correct number!")
        run = false
    }
}
```

`string.toInt()` casts string to integer. Since we already checked whether it is an integer, we can use this method.

That is our full code!

```rust
use "random"

number = random.intRange(100) + 1

run = true
while run {
    guess = input("Guess a number between 1 and 100: ")

    if !guess.isInt() {
        println("Your input is not an integer!")
        next
    }

    guess = guess.toInt()

    if guess > number
        println("Guess smaller!")
    else if guess < number
        println("Guess larger!")
    else {
        println("You guess the correct number!")
        run = false
    }
}
```
