use std::collections::HashMap;

use crate::env::Environment;
use crate::helper::*;
use crate::obj::*;

pub(crate) fn new_integer(x: i64) -> BaseObject {
    BaseObject {
        class: Class {
            name: "Integer",
            instance_vars: Environment {
                store: HashMap::default()
            },
            super_class: None,
        },
        internal_value: ValueType { int: x },
        to_string_func: wrap(|this| unsafe { this.internal_value.int.to_string() }),
        to_inspect_func: wrap(|this| unsafe { this.internal_value.int.to_string() }),
        is_truthy_func: wrap(|this| unsafe { this.internal_value.int != 0 }),
        equal_func: wrap(
            |this, other|
                this.internal_value == other.internal_value && this.class == other.class
        ),
        call_func: not_callable(),
        parent_obj: None,
    }
}
