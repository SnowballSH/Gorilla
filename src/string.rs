#![forbid(unsafe_code)]

use crate::any::any_class;
use crate::env::Environment;
use crate::obj::ValueType::*;
use crate::obj::{not_callable, BaseObject, Class};
use inner::inner;
use std::collections::HashMap;

fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Str);
    a.to_string()
}

fn k3(this: BaseObject) -> bool {
    inner!(this.internal_value, if Str) != ""
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

pub fn new_string<'a>(x: String) -> BaseObject<'a> {
    let mut _env = HashMap::default();

    BaseObject {
        class: Class {
            name: "String",
            instance_vars: Environment { store: _env },
            super_class: Some(any_class()),
        },
        internal_value: Str(x),
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}
