use lazy_static::*;
use pest::iterators::Pair;
use pest::Parser;
use pest::prec_climber::*;
use pest_derive::*;

use Rule::*;

use crate::ast::*;

lazy_static! {
    static ref PREC_CLIMBER: PrecClimber<Rule> = {
        use Assoc::*;

        PrecClimber::new(vec![
            Operator::new(add, Left) | Operator::new(sub, Left),
            Operator::new(mul, Left) | Operator::new(div, Left) | Operator::new(modulo, Left),
        ])
    };
}

#[derive(Parser)]
#[grammar = "gorilla.pest"]
pub struct GorillaParser;

fn infix<'a>(lhs: Expression<'a>, op: Pair<'a, Rule>, rhs: Expression<'a>) -> Expression<'a> {
    Expression::Infix(Box::new(Infix {
        left: lhs,
        operator: op.as_str(),
        right: rhs,
        pos: op.as_span(),
    }))
}

fn others(pair: Pair<Rule>) -> Expression {
    match pair.as_rule() {
        Rule::integer => Expression::Int(Integer {
            value: pair.as_str().parse().unwrap(),
            pos: pair.as_span(),
        }),
        Rule::identifier => Expression::GetVar(GetVar {
            name: pair.as_str(),
            pos: pair.as_span(),
        }),
        Rule::assign => {
            let mut inner = pair.clone().into_inner();
            let name = inner.next().unwrap().as_str();
            let res = inner.next().unwrap();
            Expression::SetVar(Box::new(SetVar {
                name,
                value: parse_expression(res),
                pos: pair.as_span(),
            }))
        }
        Rule::call => {
            let mut inner = pair.clone().into_inner();
            let res = inner.next().unwrap();
            let args: Vec<Pair<Rule>> = inner.collect();
            Expression::Call(
                Box::new(Call {
                    callee: parse_expression(res),
                    arguments: args
                        .into_iter().map(|w| parse_expression(w))
                        .collect(),
                    pos: pair.as_span(),
                }))
        }
        Rule::expression => climb(pair),
        _ => {
            dbg!(pair.as_rule());
            unreachable!()
        }
    }
}

pub fn climb(pair: Pair<Rule>) -> Expression {
    //dbg!(&pair);
    PREC_CLIMBER.climb(pair.into_inner(), others, infix)
}

fn parse_expression(pair: Pair<Rule>) -> Expression {
    let inner: Vec<Pair<Rule>> = pair.clone().into_inner().collect();
    let res = if inner.len() == 0 {
        others(pair)
    } else {
        climb(pair)
    };

    res
}

fn parse_statement(pair: Pair<Rule>) -> Statement {
    match pair.as_rule() {
        Rule::expression_stmt => {
            let p = pair.into_inner().next().unwrap();
            let s = p.clone().as_span();
            Statement::ExprStmt(
                ExprStmt {
                    expr: parse_expression(p),
                    pos: s,
                }
            )
        }
        _ => unreachable!()
    }
}

pub fn parse(code: &str) -> Result<Program, pest::error::Error<Rule>> {
    let res = GorillaParser::parse(Rule::program, code);
    match res {
        Ok(res) => {
            let mut ast = vec![];
            for pair in res {
                match pair.as_rule() {
                    Rule::stmt | Rule::expression_stmt => {
                        ast.push(parse_statement(pair))
                    }
                    _ => {}
                }
            }
            Ok(ast)
        }
        Err(e) => Err(e)
    }
}

#[cfg(test)]
mod tests {
    use crate::parser::parse;

    #[test]
    fn parsing() {
        let res = parse("
(println)(a + 99 % 3)");
        match res {
            Ok(x) => {
                // dbg!(&x);
                assert_eq!(x.len(), 1);
            }
            Err(x) => { panic!("{}", x.to_string()); }
        };

        let res = parse("9223372036854775808");
        match res {
            Ok(_) => {
                panic!("no error");
            }
            Err(_) => {}
        };
    }
}
