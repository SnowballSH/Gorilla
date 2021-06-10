use pest::Span;

#[derive(Debug, Clone)]
pub enum Expression<'a> {
    Int(Integer<'a>),
    String(String<'a>),
    GetVar(GetVar<'a>),
    SetVar(Box<SetVar<'a>>),
    Infix(Box<Infix<'a>>),
    Prefix(Box<Prefix<'a>>),
    Call(Box<Call<'a>>),
    GetInstance(Box<GetInstance<'a>>),
    Index(Box<Index<'a>>),
    Vec_(Box<VectorDec<'a>>),

    If(Box<If<'a>>),
    While(Box<While<'a>>),
}

#[derive(Debug, Clone)]
pub struct If<'a> {
    pub cond: Expression<'a>,
    pub body: Program<'a>,
    pub other: Program<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct While<'a> {
    pub cond: Expression<'a>,
    pub body: Program<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Infix<'a> {
    pub left: Expression<'a>,
    pub operator: &'a str,
    pub right: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Prefix<'a> {
    pub operator: &'a str,
    pub right: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Integer<'a> {
    pub value: u64,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct String<'a> {
    pub value: std::string::String,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct GetVar<'a> {
    pub name: &'a str,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct SetVar<'a> {
    pub name: &'a str,
    pub value: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Call<'a> {
    pub callee: Expression<'a>,
    pub arguments: Vec<Expression<'a>>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct Index<'a> {
    pub callee: Expression<'a>,
    pub index: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct VectorDec<'a> {
    pub values: Vec<Expression<'a>>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct GetInstance<'a> {
    pub parent: Expression<'a>,
    pub name: &'a str,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub enum Statement<'a> {
    ExprStmt(ExprStmt<'a>),
    FunctionDeclare(FunctionDeclare<'a>),
}

#[derive(Debug, Clone)]
pub struct ExprStmt<'a> {
    pub expr: Expression<'a>,
    pub pos: Span<'a>,
}

#[derive(Debug, Clone)]
pub struct FunctionDeclare<'a> {
    pub name: &'a str,
    pub args: Vec<&'a str>,
    pub body: Program<'a>,
    pub pos: Span<'a>,
}

pub type Program<'a> = Vec<Statement<'a>>;
