#![forbid(unsafe_code)]

use crate::env::Environment;
use crate::obj::ValueType::*;
use crate::obj::*;
use inner::inner;

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

pub fn new_native_function(x: NativeFunctionType) -> BaseObject {
    BaseObject {
        class: Class {
            name: "Native Function",
            instance_vars: Environment::default(),
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

#[cfg(test)]
mod tests {
    use crate::builtin_types::integer::new_integer;
    use crate::builtin_types::native_function::new_native_function;
    use crate::obj::*;

    #[test]
    fn basic() {
        fn idk<'a>(_this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
            Ok(new_integer(1))
        }
        let f = new_native_function(("idk", idk));
        assert_eq!(f.call(f.clone(), vec![]).unwrap(), new_integer(1))
    }
}
