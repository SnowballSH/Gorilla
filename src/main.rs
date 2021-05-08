#![forbid(unsafe_code)]

use std::env::args;
use std::fs::File;
use std::io::Read;
use std::io;
use crate::helpers::{run_code_with_env};
use crate::env::Environment;

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
pub mod helpers;
mod overall_test;

fn get_input() -> String{
    let mut input = String::new();
    match io::stdin().read_line(&mut input) {
        Ok(_goes_into_input_above) => {},
        Err(_no_updates_is_fine) => {},
    }
    input.trim_end().to_string()
}

fn main() {
    let argv: Vec<String> = args().collect();

    if argv.len() < 2 {
        let mut environment = Environment::default();
        loop {
            let ip = get_input();
            if ip.trim() == ":quit" {
                break;
            }
            let result = run_code_with_env(&ip, environment);
            environment = result.1;
            match result.0 {
                Ok(x) => {
                    if let Some(val) = x {
                        println!("#>> {}", val.to_inspect_string());
                    }
                }
                Err(e) => {
                    eprintln!("{}", e);
                }
            }
        }
        return;
    }

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
