use std::io::Cursor;

use crate::env::Environment;
use crate::grammar::Grammar;
use crate::integer::new_integer;
use crate::obj::*;

pub struct VM<'a> {
    pub source: Vec<u8>,
    pub ip: usize,
    pub line: usize,
    pub stack: Vec<BaseObject<'a>>,
    pub last_popped: Option<BaseObject<'a>>,
    pub env: Environment<'a>,
}

impl<'a> VM<'a> {
    pub fn new(source: Vec<u8>) -> Self {
        VM {
            source,
            ip: 0,
            line: 0,
            stack: vec![],
            last_popped: None,
            env: Default::default(),
        }
    }

    fn push(&mut self, obj: BaseObject<'a>) {
        self.stack.push(obj)
    }

    fn pop(&mut self) -> BaseObject<'a> {
        let popped = self.stack.pop().expect("Pop on empty stack...");
        self.last_popped = Some(popped.clone());
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

    fn read_string(&mut self) -> String {
        let length = self.read();
        let mut bytes = vec![];
        for _ in 0..length {
            bytes.push(self.read());
        }
        String::from_utf8(bytes).unwrap()
    }

    pub fn run(&mut self) -> Result<Option<BaseObject<'a>>, String> {
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

        Ok(self.last_popped.clone())
    }

    fn run_statement(&mut self) -> Option<String> {
        let type_ = Grammar::from(self.read());
        match type_ {
            Grammar::Advance => self.line += 1,
            Grammar::Back => self.line -= 1,
            Grammar::Pop => { self.pop(); }
            Grammar::Noop => {}
            Grammar::Integer => {
                let i = self.read_int();
                self.push(new_integer(i))
            }
            Grammar::Getvar => {
                let name = self.read_string();
                let res = self.env.get(name.clone());
                match res {
                    Some(x) => self.push(x.clone()),
                    None => return Some(format!("Variable '{}' is not defined", name))
                }
            }
            Grammar::Setvar => {
                let name = self.read_string();
                let val = self.pop();

                self.env.set(name, val.clone());
                self.push(val);
            }
            Grammar::GetInstance => {
                let self_ = self.pop();
                let g = self.read_string();
                let res = self_.instance_get(g.clone());
                match res {
                    Some(mut x) => {
                        x.set_parent(self_);
                        self.push(x);
                    }
                    None => return Some(format!(
                        "Attribute '{}' does not exist on '{}' ({})",
                        g, self_.to_string(), self_.class.to_string()
                    ))
                }
            }
            Grammar::Call => {
                let amount = self.read_int();
                let o = self.pop();
                let mut args = vec![];
                for _ in 0..amount {
                    args.push(self.pop());
                }

                let res = o.call(o.clone(), args);
                match res {
                    Err(x) => return Some(x),
                    Ok(x) => self.push(x)
                };
            }
            _ => return Some(format!("Invalid instruction: {}", type_ as u8))
        };

        None
    }
}

#[cfg(test)]
mod tests {
    use crate::grammar::Grammar;
    use crate::obj::ValueType;
    use crate::vm::VM;

    #[test]
    fn test_var() {
        let mut vm = VM::new(vec![Grammar::Magic as u8,
                                  Grammar::Integer as u8, 1, 0x03,
                                  Grammar::Setvar as u8, 2, 'a' as u8, 'b' as u8,
                                  Grammar::Pop as u8
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert!(
                x.expect("No popped").internal_value == ValueType { int: 3 }
            )
        }

        let mut vm = VM::new(vec![Grammar::Magic as u8,
                                  Grammar::Integer as u8, 1, 0x04,
                                  Grammar::Setvar as u8, 2, 'a' as u8, 'b' as u8,
                                  Grammar::Pop as u8,
                                  Grammar::Getvar as u8, 2, 'a' as u8, 'b' as u8,
                                  Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert!(
                x.expect("No popped").internal_value == ValueType { int: 4 }
            )
        }

        let mut vm = VM::new(vec![Grammar::Magic as u8,
                                  Grammar::Integer as u8, 1, 0x04,
                                  Grammar::Getvar as u8, 2, 'a' as u8, 'b' as u8,
                                  Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error...")
        }
    }

    #[test]
    fn test_call() {
        let mut vm = VM::new(vec![Grammar::Magic as u8,
                                  Grammar::Integer as u8, 1, 0x04,
                                  Grammar::Call as u8, 1, 0x00,
                                  Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error...")
        }
    }

    #[test]
    fn test_attr() {
        let mut vm = VM::new(vec![Grammar::Magic as u8,
                                  Grammar::Integer as u8, 1, 0x04,
                                  Grammar::GetInstance as u8, 1, '?' as u8,
                                  Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error...")
        }
    }
}