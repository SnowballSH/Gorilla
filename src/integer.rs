use std::collections::HashMap;

use crate::env::Environment;
use crate::obj::*;

pub(crate) fn new_integer(x: i64) -> BaseObject {
    fn k1(this: &BaseObject) -> String {
        unsafe { this.internal_value.int.to_string() }
    }
    fn k3(this: &BaseObject) -> bool {
        unsafe { this.internal_value.int != 0 }
    }
    fn k4(this: &BaseObject, other: &BaseObject) -> bool {
        this.internal_value == other.internal_value && this.class == other.class
    }
    BaseObject {
        class: Class {
            name: "Integer",
            instance_vars: Environment {
                store: HashMap::default()
            },
            super_class: None,
        },
        internal_value: ValueType { int: x },
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable(),
        parent_obj: None,
    }
}
