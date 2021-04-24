use std::io::Cursor;

use crate::env::Environment;
use crate::grammar::Grammar;
use crate::obj::*;
use crate::integer::new_integer;

pub(crate) struct VM<'a> {
    pub(crate) source: Vec<u8>,
    pub(crate) ip: usize,
    pub(crate) line: usize,
    pub(crate) stack: Vec<&'a BaseObject>,
    pub(crate) error: Option<String>,
    pub(crate) last_popped: Option<&'a BaseObject>,
    pub(crate) env: Environment,
}

impl<'a> VM<'a> {
    pub(crate) fn new(source: Vec<u8>) -> Self {
        VM {
            source,
            ip: 0,
            line: 0,
            stack: vec![],
            error: None,
            last_popped: None,
            env: Default::default(),
        }
    }

    fn push(&mut self, obj: &'a BaseObject) {
        self.stack.push(obj)
    }

    fn pop(&mut self) -> &'a BaseObject {
        let popped = self.stack.pop().expect("Pop on empty stack...");
        self.last_popped = Some(popped);
        popped
    }

    fn read(&mut self) -> u8 {
        let k = self.source[self.ip];
        self.ip += 1;
        k
    }

    fn read_int(&mut self) -> i64 {
        let length = self.read();
        let mut number = vec![];
        for _ in 0..length {
            number.push(self.read());
        }
        leb128::read::signed(&mut Cursor::new(number)).expect("Not a valid integer")
    }

    pub(crate) fn run(&mut self) -> Result<Option<&'a BaseObject>, String> {
        let length = self.source.len();
        if length == 0 || self.read() != Grammar::Magic.into() {
            return Err("Not a valid Gorilla bytecode".parse().unwrap());
        }

        while self.ip < length {
            let res = self.run_statement();
            match res {
                Some(x) => return Err(x),
                None => ()
            };
        }

        Ok(self.last_popped)
    }

    fn run_statement(&mut self) -> Option<String> {
        let type_ = Grammar::from(self.read());
        match type_ {
            Grammar::Advance => self.line += 1,
            Grammar::Back => self.line -= 1,
            Grammar::Pop => { self.pop(); },
            Grammar::Noop => {},
            Grammar::Integer => {
                let i = self.read_int();
                self.push(&new_integer(i))
            }
            _ => return Some(format!("Invalid instruction: {}", type_ as u8))
        };

        None
    }
}
