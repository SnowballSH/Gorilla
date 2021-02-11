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
    do.something1(arg1 + arg2 + arg3)
    do.something2(arg1 * arg2 * arg3)
}
func NoArgument()
	do.something()
```

### Function Literal

Syntax:

```ruby
functionLiteral ::= "fn" parameters? block
```

Syntax in code:

```rust
fn(arg1, arg2, arg3) {
    do.something1(arg1 + arg2 + arg3)
    do.something2(arg1 * arg2 * arg3)
}
fn do.something
```

#### `Function Statement` is a statement while `Function Literal` is an expression

`func name(){}` is the same as `name = fn(){}`
