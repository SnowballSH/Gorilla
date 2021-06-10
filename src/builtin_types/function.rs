use inner::inner;

use crate::builtin_types::any::any_class;
use crate::builtin_types::null::new_null;
use crate::env::Environment;
use crate::grammar::Grammar;
use crate::obj::ValueType::*;
use crate::obj::*;
use crate::vm::VM;

fn k1(this: BaseObject) -> String {
    let x = inner!(this.internal_value, if Function);
    format!("Function '{}'", x.0)
}

fn k3(_: BaseObject) -> bool {
    true
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

fn call<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, e: Environment<'a>) -> ObjResult<'a> {
    let x = inner!(this.clone().internal_value, if Function);

    if x.1.len() != args.len() {
        return Err(format!(
            "{} expected {} arguments, got {}",
            this.to_inspect_string(),
            x.1.len(),
            args.len()
        ));
    }

    let mut s = x.2;
    s.insert(0, Grammar::Magic as u8);
    let mut v = VM::new(s);
    v.env = e;

    let mut i = 0;
    for name in x.1 {
        v.env.set(name, args[i].clone());
        i += 1;
    }

    let res = v.run();

    match res {
        Ok(x) => Ok(x.unwrap_or(new_null())),
        Err(e) => Err(format!("In {}: {}", this.to_inspect_string(), e)),
    }
}

pub fn new_function<'a>(x: FunctionType) -> BaseObject<'a> {
    BaseObject {
        class: Class {
            name: "Function",
            instance_vars: Environment::default(),
            super_class: Some(any_class()),
        },
        internal_value: Function(x),
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: call,
        parent_obj: None,
    }
}
