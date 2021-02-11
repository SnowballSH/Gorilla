## Special Syntax

### Gorilla has some special, flexible syntax that you need to know

---

### Block

```ruby
block ::= "{" statement* "}"
      |   statement
```

which means it is valid to not have `{}` if you only have **one** statement.

```go
func abc() {
    a()
    b()
    c()
}

func doA()
    a()

func doB() b()

if 1 == 2
    println("NO WAY!")
```

---

### Function calls

```ruby
functionCall ::= expression "(" expression? ("," expression)* ")"
             |   doExpression
             |   functionLiteral
```

when the function only has 1 parameter, and you want to input an argument, you don't have to write `()` around the `do`
or `fn` keywords.

```ruby
abc(1, 2, 3)

[1, 2, 3].each(do(i){println(i ** 2)}) # This is ugly
```

a better-looking code with this special syntax:

```ruby
[1, 2, 3].each
    do(i)
        println(i ** 2)
# A lot better
```

[**back**](https://snowballsh.github.io/Gorilla/api)
