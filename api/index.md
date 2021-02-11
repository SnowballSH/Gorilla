## Gorilla language specification

### Top-level syntax

```ruby
terminator   ::= (";" | "\n" | ("\r" "\n") | "\r")
comment      ::= "#" [^\n\r]*

program      ::= (statement (terminator | EOF))*
```

### List of contents

- [Special Syntax Fundamentals](https://snowballsh.github.io/Gorilla/api/special) (Read this first)
- [Functions](https://snowballsh.github.io/Gorilla/api/functions)
