## Gorilla tutorial

#### Welcome to Gorilla! Let's get started!

### Hello, world!

To print hello world in Gorilla:
```javascript
display("Hello, world!")
```

Function-oriented hello world:
```go
func hello(name) {
    display("Hello, " + name + "!")
}
hello("world")
```

### Literals

1. Integers

Examples:
```ruby
1
(24)
1234567
```

The size of integer in Gorilla is Golang's int64 (64-bit int).

2. Strings

Examples:
```ruby
"a string"

"But it can be
multiple
lines"
```

The size of string in Gorilla can be as large as your pc's available memory.

3. Function Literal

Examples:
```rust
fn(param1, param2) {
    return param1 + param2
}
```

Lambda function in Gorilla has the syntax similar to Rust:
```rust
fn(param1, param2, param3, ...) { <statements> }
```
