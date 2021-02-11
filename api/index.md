## Gorilla language specification

### Top-level syntax
```bnf
terminator   ::= (";" | "\n" | ("\r" "\n") | "\r")
comment      ::= "#" [^\n\r]*
program      ::= (statement (terminator | EOF))*
```
a