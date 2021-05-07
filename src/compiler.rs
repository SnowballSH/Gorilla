use crate::ast::*;
use crate::grammar::Grammar;
use crate::helpers::*;

#[derive(Debug)]
#[doc = "The compiler struct. Compiles AST to bytecodes."]
pub struct Compiler<'a> {
    #[doc = "resulting bytecode"]
    pub result: Vec<u8>,
    #[doc = "last line processed"]
    pub last_line: usize,
    #[doc = "source code"]
    pub source: &'a str,
}

impl<'a> Compiler<'a> {
    #[doc = "Creates a new compiler"]
    pub fn new(source: &'a str) -> Self {
        Compiler {
            result: vec![Grammar::Magic as u8],
            last_line: 0,
            source,
        }
    }

    #[inline]
    fn update_line(&mut self, line: usize) {
        while line > self.last_line {
            self.last_line += 1;
            self.emit_grammar(Grammar::Advance);
        }
        while line < self.last_line {
            self.last_line -= 1;
            self.emit_grammar(Grammar::Back);
        }
    }

    #[inline]
    fn emit(&mut self, value: u8) {
        self.result.push(value)
    }

    #[inline]
    fn emit_grammar(&mut self, value: Grammar) {
        self.emit(value as u8)
    }

    #[inline]
    fn emit_all(&mut self, value: &[u8]) {
        self.result.extend(value)
    }

    #[inline]
    fn emit_unsigned_int(&mut self, value: u64) {
        let buf = leb128_unsigned(value);
        self.emit(buf.len() as u8);
        self.emit_all(&buf);
    }

    #[inline]
    fn emit_string(&mut self, value: &str) {
        self.emit_unsigned_int(value.len() as u64);
        self.emit_all(value.as_bytes());
    }

    pub fn compile(&mut self, nodes: Program) {
        for node in nodes {
            self.compile_node(node);
        }
    }

    fn compile_node(&mut self, node: Statement) {
        match node {
            Statement::ExprStmt(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.compile_expr(x.expr);
                self.emit_grammar(Grammar::Pop);
            }
        }
    }

    fn compile_expr(&mut self, expr: Expression) {
        match expr {
            Expression::Int(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.emit_grammar(Grammar::Integer);
                self.emit_unsigned_int(x.value);
            }

            Expression::GetVar(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.emit_grammar(Grammar::Getvar);
                self.emit_string(x.name);
            }
            Expression::SetVar(x) => {
                let span = span_to_line(self.source, x.pos);
                self.update_line(span);
                self.compile_expr(x.value);
                self.update_line(span);
                self.emit_grammar(Grammar::Setvar);
                self.emit_string(x.name);
            }

            Expression::Infix(x) => {
                let line = span_to_line(self.source, x.pos);
                self.update_line(line);

                self.compile_expr(x.right);
                self.update_line(line);
                self.compile_expr(x.left);
                self.update_line(line);

                self.emit_grammar(Grammar::GetInstance);
                self.emit_string(x.operator);

                self.emit_grammar(Grammar::Call);
                self.emit_unsigned_int(1);
            }

            Expression::Call(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                let mut args = x.arguments.clone();
                args.reverse();
                for n in args {
                    self.compile_expr(n);
                }
                self.compile_expr(x.callee);
                self.emit_grammar(Grammar::Call);
                self.emit_unsigned_int(x.arguments.len() as u64)
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::compiler::Compiler;
    use crate::grammar::Grammar;
    use crate::parser::parse;

    #[test]
    fn number() {
        let code = "624485";
        let mut compiler = Compiler::new(code);
        compiler.compile(parse(code).unwrap());
        assert_eq!(compiler.result, vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8, 3, 0xe5, 0x8e, 0x26,
            Grammar::Pop as u8,
        ]);
    }
}