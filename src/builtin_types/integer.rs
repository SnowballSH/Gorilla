#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;

use crate::builtin_types::any::any_class;
use crate::builtin_types::bool::new_boolean;
use crate::builtin_types::native_function::new_native_function;
use crate::env::Environment;
use crate::obj::ValueType::*;
use crate::obj::*;

#[inline]
fn add<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_integer(a + b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn sub<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_integer(a - b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn mul<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_integer(a * b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn div<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            if b == 0 {
                return Err(format!("Integer division by 0 in {} / {}", a, b));
            }
            Ok(new_integer(a / b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn mod_<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            if b == 0 {
                return Err(format!("Integer modulo by 0 in {} % {}", a, b));
            }
            Ok(new_integer(a % b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn gt<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_boolean(a > b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn lt<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_boolean(a < b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn gteq<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_boolean(a >= b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn lteq<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let other = args.first();
    match other {
        Some(x) => {
            let a = inner!(this.parent().unwrap().internal_value, if Int);
            let b = inner!(x.internal_value, if Int, else {
            let g = inner!(this.internal_value, if NativeFunction);
                return Err(format!("{} expects an integer", g.0))
            });
            Ok(new_boolean(a <= b))
        }
        None => {
            let x = inner!(this.internal_value, if NativeFunction);
            Err(format!("{} expects 1 argument, got 0", x.0))
        }
    }
}

#[inline]
fn neg<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Int);
    Ok(new_integer(-a))
}

#[inline]
fn pos<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment<'a>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Int);
    Ok(new_integer(a))
}

#[inline]
fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Int);
    a.to_string()
}

#[inline]
fn k3(this: BaseObject) -> bool {
    let a = inner!(this.internal_value, if Int);
    a != 0
}

#[inline]
fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

pub fn new_integer<'a>(x: i64) -> BaseObject<'a> {
    let mut int_env = HashMap::default();

    int_env.insert("add".to_string(), new_native_function(("Integer.+", add)));

    int_env.insert("sub".to_string(), new_native_function(("Integer.-", sub)));

    int_env.insert("mul".to_string(), new_native_function(("Integer.*", mul)));

    int_env.insert("div".to_string(), new_native_function(("Integer./", div)));

    int_env.insert(
        "modulo".to_string(),
        new_native_function(("Integer.%", mod_)),
    );

    int_env.insert("gt".to_string(), new_native_function(("Integer.>", gt)));
    int_env.insert("lt".to_string(), new_native_function(("Integer.<", lt)));
    int_env.insert(
        "gteq".to_string(),
        new_native_function(("Integer.>=", gteq)),
    );
    int_env.insert(
        "lteq".to_string(),
        new_native_function(("Integer.<=", lteq)),
    );

    int_env.insert(
        "to_neg".to_string(),
        new_native_function(("- Integer", neg)),
    );

    int_env.insert(
        "to_pos".to_string(),
        new_native_function(("+ Integer", pos)),
    );

    int_env.insert("i".to_string(), new_native_function(("Integer.i", pos)));

    BaseObject {
        class: Class {
            name: "Integer",
            instance_vars: Environment { store: int_env },
            super_class: Some(any_class()),
        },
        internal_value: Int(x),
        to_string_func: k1,
        to_inspect_func: k1,
        is_truthy_func: k3,
        equal_func: k4,
        call_func: not_callable,
        parent_obj: None,
    }
}

#[cfg(test)]
mod tests {
    use crate::builtin_types::integer::new_integer;
    use crate::env::Environment;

    #[test]
    fn binop() {
        let vv = Environment::default();

        let ii = new_integer(10);
        let mut f = ii.instance_get("add".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(1)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "11");

        let ii = new_integer(10);
        let mut f = ii.instance_get("sub".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(1)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "9");

        let ii = new_integer(10);
        let mut f = ii.instance_get("mul".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(-10)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "-100");

        let ii = new_integer(10);
        let mut f = ii.instance_get("div".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(2)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "5");

        let ii = new_integer(10);
        let mut f = ii.instance_get("modulo".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(3)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "1");

        let ii = new_integer(10);
        let mut f = ii.instance_get("div".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(0)], vv.clone());
        assert!(res.is_err());

        let ii = new_integer(10);
        let mut f = ii.instance_get("modulo".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(0)], vv.clone());
        assert!(res.is_err());

        let ii = new_integer(10);
        let mut f = ii.instance_get("modulo".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![], vv.clone());
        assert!(res.is_err()); // 0 arguments

        let ii = new_integer(10);
        let mut f = ii.instance_get("eq".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(10)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "true");

        let ii = new_integer(10);
        let mut f = ii.instance_get("neq".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(10)], vv.clone());
        assert_eq!(res.unwrap().to_string(), "false");
    }
}
