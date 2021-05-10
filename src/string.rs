#![forbid(unsafe_code)]

use crate::any::any_class;
use crate::env::Environment;
use crate::obj::ValueType::*;
use crate::obj::{not_callable, BaseObject, Class, ObjResult};
use inner::inner;
use std::collections::HashMap;
use crate::native_function::new_native_function;
use crate::integer::new_integer;

fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Str);
    a.to_string()
}

fn k2(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Str);
    let mut rs = String::new();
    for ch in a.chars() {
        rs += &*match ch {
            '\\' => r"\\".to_string(),
            '\'' => r"\'".to_string(),
            '"' => "\\\"".to_string(),
            '\n' => r"\n".to_string(),
            '\r' => r"\r".to_string(),
            '\t' => r"\t".to_string(),
            '\0' => r"\0".to_string(),
            _ => ch.to_string(),
        };
    }
    "\"".to_owned() + &*rs + "\""
}

fn k3(this: BaseObject) -> bool {
    inner!(this.internal_value, if Str) != ""
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

#[inline]
fn to_string<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Str);
    Ok(new_string(a))
}

#[inline]
fn parse_int<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Str);
    if let Ok(res) = a.parse::<i64>() {
        Ok(new_integer(res))
    } else {
        Err(format!("Cannot parse '{}' to 64-bit integer", a))
    }
}

pub fn new_string<'a>(x: String) -> BaseObject<'a> {
    let mut _env = HashMap::default();

    _env.insert("s".to_string(), new_native_function(("String.s", to_string)));
    _env.insert("i".to_string(), new_native_function(("String.i", parse_int)));

    BaseObject {
        class: Class {
            name: "String",
            instance_vars: Environment { store: _env },
            super_class: Some(any_class()),
        },
        internal_value: Str(x),
        to_string_func: k1,
        to_inspect_func: k2,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}
