#![forbid(unsafe_code)]

use std::collections::HashMap;

use crate::env::Environment;
use crate::obj::*;
use crate::obj::ValueType::*;
use inner::inner;

pub fn new_native_function<'a>(x: NativeFunctionType<'a>) -> BaseObject<'a> {
    fn k1(this: BaseObject) -> String {
        let x = inner!(this.internal_value, if NativeFunction);
        format!("Native Function {}", x.0)
    }
    fn k3(_: BaseObject) -> bool {
        true
    }
    fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
        this.internal_value == other.internal_value && this.class == other.class
    }
    fn call<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
        let x = inner!(this.internal_value, if NativeFunction);
        x.1(this.clone(), args.clone())
    }
    BaseObject {
        class: Class {
            name: "Native Function",
            instance_vars: Environment {
                store: HashMap::default()
            },
            super_class: None,
        },
        internal_value: NativeFunction(x),
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: call,
        parent_obj: None,
    }
}
