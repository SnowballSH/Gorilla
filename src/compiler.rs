use crate::grammar::Grammar;

#[derive(Default, Debug)]
pub struct Compiler {
    pub result: Vec<u8>,
    pub last_line: usize,
}

impl Compiler {
    fn update_line(&mut self, line: usize) {
        while line > self.last_line {
            self.last_line += 1;
            self.emit(Grammar::Advance as u8);
        }
        while line < self.last_line {
            self.last_line -= 1;
            self.emit(Grammar::Back as u8);
        }
    }

    fn emit(&mut self, value: u8) {
        self.result.push(value)
    }

    fn emit_string(&mut self, value: &str) {
        self.emit(value.len() as u8);

    }
}
