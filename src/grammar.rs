use strum::IntoEnumIterator;
use strum_macros::EnumIter;

use crate::grammar::Grammar::*;

#[derive(Debug, EnumIter, Copy, Clone)]
#[repr(u8)]
#[doc = "Bytecode grammar for Gorilla, ported from Go"]
pub enum Grammar {
    Magic = 0x69,
    Pop = 1,
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
    JumpTo,
    Function,
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
