## Learn Gorilla by Examples

<br>

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

---

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

<br>

---

<br>

### Simple HTTP server

Gorilla Version: 0.4.0+

Must be on local system!

It is **extremely simple** to create a HTTP server in Gorilla:

```rust
use "http"

http.get(
    "/",
    fn(writer, request) writer.write("Hello, world!")
)

http.serve(":3000")
```

That's it! Done!

After running the code, go to [localhost:3000/](http://localhost:3000/) to see the Hello World greeting from your Gorilla HTTP server!

<br>

A more complex example:

```rust
use "http"

http.get(
    "/",
    fn(writer, request) writer.write("Hello, world!")
)

let gorillaString = `<div style="text-align: center"><button onclick="alert('You clicked me!')"><h1>
Hello, Gorilla!</h1></button></div>`

http.get(
    "/gorilla",
    fn(writer, request) writer.write(gorillaString)
)

http.get(
    "/headers",
    fn(writer, request)
        request.header.each do(key, value)
            writer.write("$ " + key + " = " + value + "\n")
)

println("Listening on http://localhost:3000/")

http.serve(":3000")
```

The above code creates 3 pages: `/`, `/gorilla`, and `/headers`

`/` is the same as before

`/gorilla` shows a kind-of-beautiful button, showing that doing CSS in Gorilla is possible

`/headers` prettily prints headers from the request

See the full example [here](https://github.com/SnowballSH/Gorilla/blob/master/examples/HelloHTTP.gor)
