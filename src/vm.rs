#![forbid(unsafe_code)]

use std::io::Cursor;

use crate::bool::new_boolean;
use crate::env::Environment;
use crate::grammar::Grammar;
use crate::integer::new_integer;
use crate::obj::*;
use crate::string::new_string;

#[doc = "The Virtual Machine"]
pub struct VM<'a> {
    #[doc = "Source bytecode"]
    pub source: Vec<u8>,
    #[doc = "Input pointer"]
    pub ip: usize,
    #[doc = "Line"]
    pub line: usize,
    #[doc = "Stack of objects"]
    pub stack: Vec<BaseObject<'a>>,
    #[doc = "Last item popped. None when nothing is popped"]
    pub last_popped: Option<BaseObject<'a>>,
    #[doc = "The environment of VM"]
    pub env: Environment<'a>,
    #[doc = "Global Env"]
    pub global: Environment<'a>,
}

impl<'a> VM<'a> {
    #[doc = "New VM from vector of bytes"]
    pub fn new(source: Vec<u8>) -> Self {
        let mut global = Environment::default();
        global.set("true".to_string(), new_boolean(true));
        global.set("false".to_string(), new_boolean(false));

        VM {
            source,
            ip: 0,
            line: 0,
            stack: vec![],
            last_popped: None,
            env: Default::default(),
            global,
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

    #[doc = "Run the bytecode"]
    pub fn run(&mut self) -> Result<Option<BaseObject<'a>>, String> {
        let length = self.source.len();
        if length == 0 || self.read() != Grammar::Magic.into() {
            return Err("Not a valid Gorilla bytecode".parse().unwrap());
        }

        while self.ip < length {
            let res = self.run_statement();
            match res {
                Some(x) => return Err(x),
                None => (),
            };
        }

        Ok(self.last_popped.clone())
    }

    fn run_statement(&mut self) -> Option<String> {
        let type_ = Grammar::from(self.read());
        match type_ {
            Grammar::Advance => self.line += 1,
            Grammar::Back => self.line -= 1,
            Grammar::Pop => {
                self.pop();
            }
            Grammar::Noop => {}

            Grammar::Integer => {
                let i = self.read_int();
                self.push(new_integer(i));
            }

            Grammar::String => {
                let i = self.read_string();
                self.push(new_string(i));
            }

            Grammar::Getvar => {
                let name = self.read_string();
                let res = self.env.get(name.clone());
                match res {
                    Some(x) => self.push(x.clone()),
                    None => {
                        let res = self.global.get(name.clone());
                        match res {
                            Some(x) => self.push(x.clone()),
                            None => return Some(format!("Variable '{}' is not defined", name)),
                        }
                    }
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
                    None => {
                        return Some(format!(
                            "Attribute '{}' does not exist on '{}' ({})",
                            g,
                            self_.to_string(),
                            self_.class.to_string()
                        ))
                    }
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
                    Ok(x) => self.push(x),
                };
            }
            Grammar::Jump => {
                let where_ = self.read_int();
                self.ip = (where_ + 1) as usize
            }
            Grammar::JumpIfFalse => {
                let where_ = self.read_int();
                if !self.pop().is_truthy() {
                    self.ip = (where_ + 1) as usize
                }
            }
            _ => return Some(format!("Invalid instruction: {}", type_ as u8)),
        };

        None
    }
}

#[cfg(test)]
mod tests {
    use crate::grammar::Grammar;
    use crate::obj::ValueType::*;
    use crate::vm::VM;

    #[test]
    fn test_var() {
        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x03,
            Grammar::Setvar as u8,
            2,
            'a' as u8,
            'b' as u8,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert_eq!(x.expect("No popped").internal_value, Int(3)),
        }

        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x04,
            Grammar::Setvar as u8,
            2,
            'a' as u8,
            'b' as u8,
            Grammar::Pop as u8,
            Grammar::Getvar as u8,
            2,
            'a' as u8,
            'b' as u8,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert_eq!(x.expect("No popped").internal_value, Int(4)),
        }

        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x04,
            Grammar::Getvar as u8,
            2,
            'a' as u8,
            'b' as u8,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error..."),
        }
    }

    #[test]
    fn test_call() {
        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x04,
            Grammar::Call as u8,
            1,
            0x00,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error..."),
        }
    }

    #[test]
    fn test_attr() {
        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x04,
            Grammar::GetInstance as u8,
            1,
            '?' as u8,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(_) => {}
            Ok(_) => panic!("No Error..."),
        }
    }

    #[test]
    fn test_arth() {
        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x03,
            Grammar::Integer as u8,
            1,
            0x06,
            Grammar::GetInstance as u8,
            1,
            '+' as u8,
            Grammar::Call as u8,
            1,
            0x01,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert_eq!(x.expect("No popped").internal_value, Int(9)),
        }
    }

    #[test]
    fn test_jump() {
        let mut vm = VM::new(vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8,
            1,
            0x01,
            Grammar::JumpIfFalse as u8,
            1,
            12,
            Grammar::Integer as u8,
            1,
            0x03,
            Grammar::Jump as u8,
            1,
            20,
            Grammar::Advance as u8,
            Grammar::Advance as u8,
            Grammar::Integer as u8,
            1,
            0x04,
            Grammar::Noop as u8,
            Grammar::Back as u8,
            Grammar::Back as u8,
            Grammar::Pop as u8,
        ]);
        let res = vm.run();
        match res {
            Err(x) => panic!("Error: {}", x),
            Ok(x) => assert_eq!(x.expect("No popped").internal_value, Int(3)),
        }
    }
}
