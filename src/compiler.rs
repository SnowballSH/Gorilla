use crate::ast::*;
use crate::grammar::Grammar;
use crate::helpers::span_to_line;

#[derive(Debug)]
pub struct Compiler<'a> {
    pub result: Vec<u8>,
    pub last_line: usize,
    pub source: &'a str,
}

impl<'a> Compiler<'a> {
    fn new(source: &'a str) -> Self {
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
    fn emit_int(&mut self, value: i64) {
        let mut buf = [0; 1024];
        let mut writable = &mut buf[..];
        leb128::write::signed(&mut writable, value).expect("Should write number");
        self.emit_all(&buf);
    }

    #[inline]
    fn emit_unsigned_int(&mut self, value: u64) {
        let mut buf = [0; 1024];
        let mut writable = &mut buf[..];
        leb128::write::unsigned(&mut writable, value).expect("Should write number");
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
                self.emit_int(x.value);
            }
            Expression::GetVar(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.emit_grammar(Grammar::Getvar);
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
