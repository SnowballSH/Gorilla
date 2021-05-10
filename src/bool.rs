#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;
use lazy_static::*;

use crate::any::any_class;
use crate::env::Environment;
use crate::integer::new_integer;
use crate::obj::{BaseObject, Class, not_callable, ObjResult};
use crate::obj::ValueType::*;
use crate::native_function::new_native_function;

fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Bool);
    a.to_string()
}

fn k3(this: BaseObject) -> bool {
    inner!(this.internal_value, if Bool)
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

#[inline]
fn to_int<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Bool);
    Ok(new_integer(if a { 1 } else { 0 }))
}

pub fn new_boolean<'a>(x: bool) -> BaseObject<'a> {
    let mut _env = HashMap::default();

    _env.insert("i".to_string(), new_native_function(("Boolean.to_i", to_int)));

    BaseObject {
        class: Class {
            name: "Boolean",
            instance_vars: Environment { store: _env },
            super_class: Some(any_class()),
        },
        internal_value: Bool(x),
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}

lazy_static! {
    pub static ref GORILLA_TRUE: BaseObject<'static> = new_boolean(true);
    pub static ref GORILLA_FALSE: BaseObject<'static> = new_boolean(false);
}

#[cfg(test)]
mod tests {
    use crate::bool::*;

    #[test]
    fn test() {
        assert_eq!(GORILLA_TRUE.to_string(), "true");
        assert_eq!(GORILLA_FALSE.to_string(), "false");
    }
}
