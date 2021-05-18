#![forbid(unsafe_code)]

use std::collections::HashMap;

use crate::any::any_class;
use crate::env::Environment;
use crate::obj::{BaseObject, Class, not_callable, ObjResult};
use crate::obj::ValueType::*;
use crate::native_function::new_native_function;
use crate::string::new_string;

fn k1(_: BaseObject) -> String {
    "null".to_string()
}

fn k3(_: BaseObject) -> bool {
    false
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

#[inline]
fn to_string<'a>(_: BaseObject<'a>, _: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    Ok(new_string("null".to_string()))
}

pub fn new_null<'a>() -> BaseObject<'a> {
    let mut _env = HashMap::default();

    _env.insert("s".to_string(), new_native_function(("Null.s", to_string)));

    BaseObject {
        class: Class {
            name: "Null",
            instance_vars: Environment { store: _env },
            super_class: Some(any_class()),
        },
        internal_value: Null,
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}
