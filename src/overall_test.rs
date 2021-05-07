#[cfg(test)]
mod test_overall {
    use crate::parser::parse;
    use crate::compiler::Compiler;
    use crate::vm::VM;
    use crate::integer::new_integer;

    #[test]
    fn connect() {
        let code = "1453 + 26 * 7";
        let mut compiler = Compiler::new(code);
        compiler.compile(parse(code).unwrap());
        let mut vm = VM::new(compiler.result);
        let result = vm.run();
        assert_eq!(result, Ok(Some(new_integer(1453 + 26 * 7))));

        let code = "abc = 4763 / 23 + 98765; abc";
        let mut compiler = Compiler::new(code);
        compiler.compile(parse(code).unwrap());
        let mut vm = VM::new(compiler.result);
        let result = vm.run();
        assert_eq!(result, Ok(Some(new_integer(4763 / 23 + 98765))));
    }
}