use pest::Span;

#[derive(Debug, Clone)]
pub enum Expression<'a> {
    Int(Integer<'a>),
    GetVar(GetVar<'a>),
    Infix(Box<Infix<'a>>),
    Call(Box<Call<'a>>), // callee, args
}

#[derive(Debug, Clone)]
pub struct Infix<'a> {
    pub left: Expression<'a>,
    pub operator: &'a str,
    pub right: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Integer<'a> {
    pub value: i64,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct GetVar<'a> {
    pub name: &'a str,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Call<'a> {
    pub callee: Expression<'a>,
    pub arguments: Vec<Expression<'a>>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub enum Statement<'a> {
    ExprStmt(ExprStmt<'a>),
}

#[derive(Debug, Clone)]
pub struct ExprStmt<'a> {
    pub expr: Expression<'a>,
    pub pos: Span<'a>,
}

pub type Program<'a> = Vec<Statement<'a>>;