use crate::tokens::{Token, Data};

pub struct Lexer {
    input: &'static str,
    input_size: usize,

    pos: usize,
    read_pos: usize,
    ch: char,
    line_count: usize,
    char_count: usize,

    tokens: Vec<Token>,
}

impl Lexer {
    pub(crate) fn new(input: &'static str) -> Lexer {
        Lexer {
            input,
            input_size: input.chars().count(),
            pos: 0,
            read_pos: 0,
            ch: '\0',
            line_count: 0,
            char_count: 0,
            tokens: vec![],
        }
    }

    pub fn next_token(&mut self) -> Token {
        self.skip_whitespace();

        let tok = match self.ch {
            _ => Data::Illegal
        };

        self.read_char();

        return Token {
            _data: tok,
            _line: self.line_count,
            _char: self.char_count,
        };
    }

    fn is_whitespace(&self, ch: char) -> bool {
        return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r';
    }

    fn skip_whitespace(&mut self) {
        loop {
            if !self.is_whitespace(self.ch) {
                break;
            }

            self.read_char();
        }
    }

    fn read_char(&mut self) {
        if self.read_pos >= self.input_size {
            self.ch = '\0';
        } else {
            self.ch = self.input.chars().nth(self.read_pos).unwrap();
        }

        self.pos = self.read_pos;
        self.read_pos += 1;
    }
}

// TESTS

#[test]
fn test_lexer() {
    let mut lexer = Lexer::new("23");
    let res = lexer.next_token();
    assert_eq!(res._data, Data::Illegal)
}
