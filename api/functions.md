## Functions

There are two kinds of functions in Gorilla

### Function Statement

Syntax:

```ruby
functionStatement ::= "func" identifier parameters block
```

Syntax in code:

```go
func FunctionName(arg1, arg2, arg3) {
    doo.something1(arg1 + arg2 + arg3)
    doo.something2(arg1 * arg2 * arg3)
}
func NoArgument()
    doo.something()
```

### Function Literal

Syntax:

```ruby
functionLiteral ::= "fn" parameters? block
```

Syntax in code:

```rust
fn(arg1, arg2, arg3) {
    doo.something1(arg1 + arg2 + arg3)
    doo.something2(arg1 * arg2 * arg3)
}
fn doo.something
```

#### `Function Statement` is a statement while `Function Literal` is an expression

`func name(){}` is the same as `name = fn(){}`

[**back**](https://snowballsh.github.io/Gorilla/api)
