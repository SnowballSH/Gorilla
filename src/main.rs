use std::env::args;
use std::fs::File;
use std::io::Read;

mod obj;
mod env;
mod integer;
mod helper;
mod vm;
mod grammar;

fn main() {
    let argv: Vec<String> = args().collect();
    let filename = &argv[1];
    let mut file = File::open(filename).expect("Unable to open the file");
    let mut contents = String::new();
    file.read_to_string(&mut contents).expect("Unable to read the file");

    let mut vm = vm::VM::new(Vec::from(contents));
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
