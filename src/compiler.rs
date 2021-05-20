use crate::ast::*;
use crate::grammar::Grammar;
use crate::helpers::*;

#[derive(Debug)]
/// The compiler struct. Compiles AST to bytecodes.
pub struct Compiler<'a> {
    /// resulting bytecode
    pub result: Vec<u8>,
    /// last line processed
    pub last_line: usize,
    /// source code
    pub source: &'a str,
}

impl<'a> Compiler<'a> {
    /// Creates a new compiler
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

    /// compiles a full program
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
            Statement::FunctionDeclare(x) => {
                self.update_line(span_to_line(self.source, x.pos));

                self.emit_grammar(Grammar::Function);
                self.emit_unsigned_int(x.args.len() as u64);

                for s in x.args {
                    self.emit_string(s);
                }

                let mut comp = Compiler {
                    result: vec![],
                    last_line: self.last_line,
                    source: self.source,
                };

                comp.compile(x.body);
                self.emit_unsigned_int(comp.result.len() as u64);
                self.emit_all(comp.result.as_slice());
                self.update_line(comp.last_line);
            }
        }
    }

    //noinspection DuplicatedCode
    fn compile_expr(&mut self, expr: Expression) {
        match expr {
            Expression::Int(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.emit_grammar(Grammar::Integer);
                self.emit_unsigned_int(x.value);
            }
            Expression::String(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.emit_grammar(Grammar::String);
                self.emit_string(&*x.value);
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

            Expression::Prefix(x) => {
                let line = span_to_line(self.source, x.pos);
                self.update_line(line);

                self.compile_expr(x.right);
                self.update_line(line);

                self.emit_grammar(Grammar::GetInstance);
                self.emit_string(&*(x.operator.to_owned() + "@"));

                self.emit_grammar(Grammar::Call);
                self.emit_unsigned_int(0);
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

            Expression::GetInstance(x) => {
                self.update_line(span_to_line(self.source, x.pos));
                self.compile_expr(x.parent);
                self.emit_grammar(Grammar::GetInstance);
                self.emit_string(x.name);
            }

            Expression::If(x) => {
                let p = x.pos;
                let pos = span_to_line(self.source, p);
                self.update_line(pos);
                self.compile_expr(x.cond);

                {
                    let mut compi = Compiler {
                        result: vec![],
                        last_line: self.last_line,
                        source: self.source,
                    };

                    let mut compe = Compiler {
                        result: vec![],
                        last_line: self.last_line,
                        source: self.source,
                    };

                    compi.compile(x.body);
                    compe.compile(x.other);

                    if compi.result.len() > 0 && *compi.result.last().unwrap() == Grammar::Pop as u8 {
                        compi.result.pop();
                        compi.emit_grammar(Grammar::Noop);
                    } else {
                        compi.emit_grammar(Grammar::Null)
                    }

                    if compe.result.len() > 0 && *compe.result.last().unwrap() == Grammar::Pop as u8 {
                        compe.result.pop();
                        compe.emit_grammar(Grammar::Noop);
                    } else {
                        compe.emit_grammar(Grammar::Null)
                    }

                    compi.update_line(pos);
                    compe.update_line(pos);

                    self.emit_grammar(Grammar::JumpIfFalse);

                    compi.emit_grammar(Grammar::Jump);

                    let k = (self.result.len() + compi.result.len() + 1) as u64;
                    compi.emit_unsigned_int(
                        (
                            self.result.len() + 2 +
                                leb128_unsigned(
                                    k
                                ).len() + compi.result.len() + compe.result.len()
                        ) as u64
                    );

                    let k = (self.result.len() + compi.result.len() + 1) as u64;
                    self.emit_unsigned_int(k);

                    self.result.extend(compi.result);
                    self.result.extend(compe.result);
                }
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

    #[test]
    fn if_else() {
        let code = "if 1 {
        } else {5
        }
        1";
        let mut compiler = Compiler::new(code);
        compiler.compile(parse(code).unwrap_or_else(|x| panic!("{}", x.to_string())));
        assert_eq!(compiler.result, vec![
            Grammar::Magic as u8,
            Grammar::Integer as u8, 1, 1,
            Grammar::JumpIfFalse as u8, 1, 10,
            Grammar::Null as u8,
            Grammar::Jump as u8, 1, 16,
            Grammar::Advance as u8,// Grammar::Advance as u8,
            Grammar::Integer as u8, 1, 5, Grammar::Noop as u8,
            Grammar::Back as u8,// Grammar::Back as u8,
            Grammar::Pop as u8,
            Grammar::Advance as u8, Grammar::Advance as u8,
            Grammar::Advance as u8,// Grammar::Advance as u8,
            Grammar::Integer as u8, 1, 1,
            Grammar::Pop as u8,
        ]);
    }
}