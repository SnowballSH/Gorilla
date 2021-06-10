#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;

use crate::builtin_types::any::any_class;
use crate::builtin_types::integer::new_integer;
use crate::builtin_types::native_function::new_native_function;
use crate::env::Environment;
use crate::obj::{BaseObject, Class, not_callable, ObjResult};
use crate::obj::ValueType::*;

fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Vector);
    let mut rs = String::new();
    rs += "[";
    rs += &a.iter().map(|x| x.to_string()).collect::<Vec<String>>().join(", ");
    rs += "]";
    rs
}

fn k2(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Vector);
    let mut rs = String::new();
    rs += "[";
    rs += &a.iter().map(|x| x.to_inspect_string()).collect::<Vec<String>>().join(", ");
    rs += "]";
    rs
}

fn k3(this: BaseObject) -> bool {
    inner!(this.internal_value, if Vector).len() != 0
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

#[inline]
fn add<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let mut a = inner!(this.parent().unwrap().internal_value, if Vector);
            let b = inner!(&x.internal_value, if Vector, else {
                let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects a vector", g.0))
            });
            for y in b {
                a.push(y.clone());
            }
            Ok(new_vector(a))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn get_index<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Vector);
            let b = inner!(&x.internal_value, if Int, else {
                let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            let c = if *b < 0 { a.len() as i64 + *b } else { *b };
            let res = a.get(c as usize);
            match res {
                Some(res) => {
                    Ok((*res).clone())
                }
                None => Err(format!("Vector index {} is out of range", b))
            }
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn length<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Vector);
    Ok(new_integer(a.len() as i64))
}

pub fn new_vector(x: Vec<BaseObject>) -> BaseObject {
    let mut _env = HashMap::default();
    _env.insert("len".to_string(), new_native_function(("Vec.len", length)));

    _env.insert("add".to_string(), new_native_function(("Vec.+", add)));
    _env.insert("get_index".to_string(), new_native_function(("Vec.get_index", get_index)));

    BaseObject {
        class: Class {
            name: "Vec",
            instance_vars: Environment { store: _env },
            super_class: Some(any_class()),
        },
        internal_value: Vector(x),
        to_string_func: k1,
        to_inspect_func: k2,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}