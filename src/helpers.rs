use pest::Span;

use crate::compiler::Compiler;
use crate::env::Environment;
use crate::obj::BaseObject;
use crate::parser::parse;
use crate::vm::VM;

#[inline]
pub fn span_to_line(source: &str, span: Span) -> usize {
    source[..span.start()].matches("\n").count()
}

pub fn leb128_unsigned(val: u64) -> Vec<u8> {
    let mut value = val | 0;
    if value < 0x80 {
        return vec![value as u8];
    }

    let mut res = vec![];

    loop {
        let mut c = (value & 0x7f) as u8;
        value >>= 7;
        if value != 0 {
            c |= 0x80;
        }
        res.push(c);
        if c & 0x80 == 0 {
            break;
        }
    }

    res
}

pub fn run_code(code: &str) -> Result<Option<BaseObject>, (String, usize)> {
    let mut compiler = Compiler::new(code);
    let p = parse(code);
    if let Err(e) = p {
        return Err((e.to_string(), 0));
    }
    compiler.compile(p.unwrap());
    let mut vm = VM::new(compiler.result);
    let result = vm.run();
    if let Ok(x) = result {
        Ok(x)
    } else if let Err(x) = result {
        Err((x, vm.line))
    } else {
        unreachable!()
    }
}

pub fn compile_code(code: &str) -> Result<Vec<u8>, (String, usize)> {
    let mut compiler = Compiler::new(code);
    let p = parse(code);
    if let Err(e) = p {
        return Err((e.to_string(), 0));
    }
    compiler.compile(p.unwrap());
    return Ok(compiler.result);
}

pub fn run_bytecode<'a>(code: Vec<u8>) -> Result<Option<BaseObject<'a>>, (String, usize)> {
    let mut vm = VM::new(code);
    let result = vm.run();
    if let Ok(x) = result {
        Ok(x)
    } else if let Err(x) = result {
        Err((x, vm.line))
    } else {
        unreachable!()
    }
}

pub fn run_code_with_env<'a>(code: &str, env: Environment<'a>)
                             -> (Result<Option<BaseObject<'a>>, (String, usize)>, Environment<'a>) {
    let mut compiler = Compiler::new(code);
    let p = parse(code);
    if let Err(e) = p {
        return (Err((e.to_string(), 0)), env);
    }
    compiler.compile(p.unwrap());
    let mut vm = VM::new(compiler.result);
    vm.env = env;
    let result = vm.run();
    if let Ok(x) = result {
        (Ok(x), vm.env)
    } else if let Err(x) = result {
        (Err((x, vm.line)), vm.env)
    } else {
        unreachable!()
    }
}
