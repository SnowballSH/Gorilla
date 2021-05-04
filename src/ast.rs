use pest::Span;

#[derive(Debug)]
pub enum Expression<'a> {
    Int(Integer<'a>),
    GetVar(GetVar<'a>),
    Infix(Box<Infix<'a>>),
    Call(Box<Call<'a>>), // callee, args
}

#[derive(Debug)]
pub struct Infix<'a> {
    pub left: Expression<'a>,
    pub operator: &'a str,
    pub right: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug)]
pub struct Integer<'a> {
    pub value: i64,
    pub pos: Span<'a>,
}

#[derive(Debug)]
pub struct GetVar<'a> {
    pub name: &'a str,
    pub pos: Span<'a>,
}

#[derive(Debug)]
pub struct Call<'a> {
    pub callee: Expression<'a>,
    pub arguments: Vec<Expression<'a>>,
    pub pos: Span<'a>,
}

#[derive(Debug)]
pub enum Statement<'a> {
    ExprStmt(Expression<'a>),
}

pub type Program<'a> = Vec<Statement<'a>>;