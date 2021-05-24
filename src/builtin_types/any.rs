#![forbid(unsafe_code)]

use inner::inner;

use crate::builtin_types::bool::new_boolean;
use crate::env::Environment;
use crate::builtin_types::native_function::new_native_function;
use crate::obj::*;
use crate::obj::ValueType::*;
use crate::builtin_types::string::new_string;

fn dbeq<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => Ok(new_boolean(this.parent().unwrap().equal_to(x.clone()))),
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

fn neq<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => Ok(new_boolean(!this.parent().unwrap().equal_to(x.clone()))),
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn to_string<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    Ok(new_string(this.parent().unwrap().to_string()))
}

#[inline]
fn to_inspect_string<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    Ok(new_string(this.parent().unwrap().to_inspect_string()))
}

pub fn any_class<'a>() -> Box<Class<'a>> {
    let mut store = Environment::default();
    store.set("==".to_string(), new_native_function(("Object.==", dbeq)));
    store.set("!=".to_string(), new_native_function(("Object.!=", neq)));
    store.set("s".to_string(), new_native_function(("Object.s", to_string)));
    store.set("inspect".to_string(), new_native_function(("Object.inspect", to_inspect_string)));
    Box::new(Class {
        name: "Object",
        instance_vars: store,
        super_class: None,
    })
}
