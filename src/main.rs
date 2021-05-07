#![forbid(unsafe_code)]

use std::env::args;
use std::fs::File;
use std::io::Read;

pub mod any;
pub mod bool;
pub mod env;
pub mod grammar;
pub mod integer;
pub mod string;
pub mod native_function;
pub mod obj;
pub mod vm;
pub mod parser;
pub mod ast;
pub mod compiler;
mod helpers;
mod overall_test;

fn main() {
    let argv: Vec<String> = args().collect();
    let filename = &argv[1];
    let mut file = File::open(filename).expect("Unable to open the file");
    let mut contents = vec![];
    file.read_to_end(&mut contents)
        .expect("Unable to read the file");

    let mut vm = vm::VM::new(contents);
    let res = vm.run();
    match res {
        Ok(x) => {
            if let Some(y) = x {
                println!("| Last item popped: {}", y.to_string());
            }
        }
        Err(e) => {
            println!("| Error: {}", e);
        }
    };
}
