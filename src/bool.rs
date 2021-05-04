#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;

use crate::any::any_class;
use crate::env::Environment;
use crate::obj::{BaseObject, Class, not_callable};
use crate::obj::ValueType::*;

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

pub fn new_boolean<'a>(x: bool) -> BaseObject<'a> {
    let mut _env = HashMap::default();

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
