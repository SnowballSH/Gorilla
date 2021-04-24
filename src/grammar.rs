use strum::IntoEnumIterator;
use strum_macros::EnumIter;

use crate::grammar::Grammar::Noop;

#[derive(Debug, EnumIter, Copy, Clone)]
pub enum Grammar {
    Magic = 0x69,
    Pop = 0,
    Noop,
    Advance,
    Back,
    Integer,
    String,
    Null,
    Getvar,
    Setvar,
    GetInstance,
    Call,
    JumpIfFalse,
    Jump,
    Lambda,
    Closure,
}

impl From<Grammar> for u8 {
    fn from(g: Grammar) -> Self {
        g as u8
    }
}

impl From<u8> for Grammar {
    fn from(x: u8) -> Self {
        for grm in Grammar::iter() {
            if grm as u8 == x {
                return grm;
            }
        }

        Noop
    }
}
