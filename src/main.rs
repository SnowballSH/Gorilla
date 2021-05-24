#![forbid(unsafe_code)]

use std::env::args;
use std::fs::File;
use std::io::Read;
use std::{io, thread};
use crate::helpers::{run_code_with_env, run_code};
use crate::env::Environment;
use console::style;

pub mod env;
pub mod grammar;
pub mod builtin_types;
pub mod modules;
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

fn _main() {
    let argv: Vec<String> = args().collect();

    if argv.len() < 2 {
        let mut environment = Environment::default();

        println!("{}", style("Welcome to Gorilla repl. Type :quit to quit.").yellow());

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
                        println!("=> {}", style(val.to_inspect_string()).blue());
                    }
                }
                Err(e) => {
                    println!("{}", style(e).red());
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

    let res = run_code(std::str::from_utf8(&*contents).unwrap());
    match res {
        Ok(_) => {}
        Err(e) => {
            println!("| In Line {}:\n| Error: {}", e.1 + 1, e.0);
        }
    };
}

static STACK_SIZE: usize = 1 << 24;

fn main() {
    let child = thread::Builder::new()
        .stack_size(STACK_SIZE)
        .name(format!("Gorilla Runtime"))
        .spawn(_main)
        .unwrap();

    child.join().unwrap();
}
