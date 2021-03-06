#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;
use unicode_segmentation::UnicodeSegmentation;

use crate::builtin_types::any::any_class;
use crate::builtin_types::integer::new_integer;
use crate::builtin_types::native_function::new_native_function;
use crate::builtin_types::vec::new_vector;
use crate::env::Environment;
use crate::obj::ValueType::*;
use crate::obj::{not_callable, BaseObject, Class, ObjResult};

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
fn parse_int<'a>(
    this: BaseObject<'a>,
    _args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Str);
    if let Ok(res) = a.parse::<i64>() {
        Ok(new_integer(res))
    } else {
        Err(format!("Cannot parse '{}' to 64-bit integer", a))
    }
}

#[inline]
fn add<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Str);
            let b = inner!(&x.internal_value, if Str, else {
                let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects a string", g.0))
            });
            Ok(new_string(a + &*b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn get_index<'a>(
    this: BaseObject<'a>,
    args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Str);
            let b = inner!(&x.internal_value, if Int, else {
                let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            let c = if *b < 0 {
                a.graphemes(true).count() as i64 + *b
            } else {
                *b
            };
            let res = a.chars().nth(c as usize);
            match res {
                Some(res) => Ok(new_string(res.to_string())),
                None => Err(format!("String index {} is out of range", b)),
            }
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn length<'a>(
    this: BaseObject<'a>,
    _args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Str);
    Ok(new_integer(a.graphemes(true).count() as i64))
}

#[inline]
fn chars<'a>(
    this: BaseObject<'a>,
    _args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let mut v = vec![];
    for x in inner!(this.parent().unwrap().internal_value, if Str).chars() {
        v.push(new_string(x.to_string()));
    }
    Ok(new_vector(v))
}

#[inline]
fn split<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Str);
            let b = inner!(&x.internal_value, if Str, else {
                let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects a string", g.0))
            });
            let mut v = vec![];
            for x in a.split(b) {
                v.push(new_string(x.to_string()));
            }
            Ok(new_vector(v))
        }
        None => {
            let mut v = vec![];
            for x in inner!(this.parent().unwrap().internal_value, if Str).split_ascii_whitespace()
            {
                v.push(new_string(x.to_string()));
            }
            Ok(new_vector(v))
        }
    }
}

pub fn new_string<'a>(x: String) -> BaseObject<'a> {
    let mut _env = HashMap::default();

    _env.insert(
        "i".to_string(),
        new_native_function(("String.i", parse_int)),
    );
    _env.insert(
        "len".to_string(),
        new_native_function(("String.len", length)),
    );

    _env.insert("add".to_string(), new_native_function(("String.+", add)));
    _env.insert(
        "get_index".to_string(),
        new_native_function(("String.get_index", get_index)),
    );
    _env.insert(
        "chars".to_string(),
        new_native_function(("String.chars", chars)),
    );
    _env.insert(
        "split".to_string(),
        new_native_function(("String.split", split)),
    );

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

#[cfg(test)]
mod tests {
    use crate::builtin_types::integer::new_integer;
    use crate::builtin_types::string::new_string;
    use crate::env::Environment;

    #[test]
    fn methods() {
        let vv = Environment::default();

        let my_string = new_string("Hello, 世界!".to_string());

        let mut f = my_string.instance_get("add".to_string()).unwrap();
        f.set_parent(my_string);
        let res = f
            .clone()
            .call(f, vec![new_string("?!".to_string())], vv.clone());
        assert_eq!(res.unwrap().to_string(), "Hello, 世界!?!");

        let my_string = new_string("Hello, 世界!".to_string());

        let mut f = my_string.instance_get("len".to_string()).unwrap();
        f.set_parent(my_string);
        let res = f.clone().call(f, vec![], vv.clone());
        assert_eq!(res.unwrap(), new_integer(10));

        let my_string = new_string("Hello, 世界!".to_string());

        let mut f = my_string.instance_get("get_index".to_string()).unwrap();
        f.set_parent(my_string);
        let res = f.clone().call(f.clone(), vec![new_integer(-2)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "界");
        let res = f.clone().call(f, vec![new_integer(2)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "l");

        let my_string = new_string("1230".to_string());

        let mut f = my_string.instance_get("i".to_string()).unwrap();
        f.set_parent(my_string);
        let res = f.clone().call(f, vec![], vv.clone());
        assert_eq!(res.unwrap(), new_integer(1230));
    }
}
