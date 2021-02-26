pub struct Token {
    pub(crate) _data: Data,
    pub(crate) _line: usize,
    pub(crate) _char: usize,
}

#[derive(PartialEq, Eq, Clone, Hash, Debug)]
pub enum Data {
    Integer(String),
    Iden(String),
    Add,
    Sub,
    Div,
    Illegal
}