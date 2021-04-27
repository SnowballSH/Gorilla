#![forbid(unsafe_code)]

use std::collections::HashMap;

use inner::inner;

use crate::any::any_class;
use crate::env::Environment;
use crate::native_function::new_native_function;
use crate::obj::ValueType::*;
use crate::obj::*;

fn add<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
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

fn sub<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
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

fn mul<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
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

fn div<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
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

fn mod_<'a>(this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
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

fn neg<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Int);
    Ok(new_integer(-a))
}

fn pos<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
    let a = inner!(this.parent().unwrap().internal_value, if Int);
    Ok(new_integer(a))
}

fn k1(this: BaseObject) -> String {
    let a = inner!(this.internal_value, if Int);
    a.to_string()
}

fn k3(this: BaseObject) -> bool {
    let a = inner!(this.internal_value, if Int);
    a != 0
}

fn k4<'a>(this: BaseObject<'a>, other: BaseObject<'a>) -> bool {
    this.internal_value == other.internal_value && this.class == other.class
}

pub fn new_integer<'a>(x: i64) -> BaseObject<'a> {
    let mut int_env = HashMap::default();

    int_env.insert("+".to_string(), new_native_function(("Integer.+", add)));

    int_env.insert("-".to_string(), new_native_function(("Integer.-", sub)));

    int_env.insert("*".to_string(), new_native_function(("Integer.*", mul)));

    int_env.insert("/".to_string(), new_native_function(("Integer./", div)));

    int_env.insert("%".to_string(), new_native_function(("Integer.%", mod_)));

    int_env.insert("-@".to_string(), new_native_function(("- Integer", neg)));

    int_env.insert("+@".to_string(), new_native_function(("+ Integer", pos)));

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
    use crate::integer::new_integer;

    #[test]
    fn binop() {
        let ii = new_integer(10);
        let mut f = ii.instance_get("+".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(1)]);
        assert_eq!(res.unwrap().to_string(), "11");

        let ii = new_integer(10);
        let mut f = ii.instance_get("-".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(1)]);
        assert_eq!(res.unwrap().to_string(), "9");

        let ii = new_integer(10);
        let mut f = ii.instance_get("*".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(-10)]);
        assert_eq!(res.unwrap().to_string(), "-100");

        let ii = new_integer(10);
        let mut f = ii.instance_get("/".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(2)]);
        assert_eq!(res.unwrap().to_string(), "5");

        let ii = new_integer(10);
        let mut f = ii.instance_get("%".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(3)]);
        assert_eq!(res.unwrap().to_string(), "1");

        let ii = new_integer(10);
        let mut f = ii.instance_get("/".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(0)]);
        assert!(res.is_err());

        let ii = new_integer(10);
        let mut f = ii.instance_get("%".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(0)]);
        assert!(res.is_err());

        let ii = new_integer(10);
        let mut f = ii.instance_get("%".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![]);
        assert!(res.is_err()); // 0 arguments

        let ii = new_integer(10);
        let mut f = ii.instance_get("==".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(10)]);
        assert_eq!(res.unwrap().to_string(), "true");

        let ii = new_integer(10);
        let mut f = ii.instance_get("!=".to_string()).unwrap();
        f.set_parent(ii);
        let res = f.clone().call(f, vec![new_integer(10)]);
        assert_eq!(res.unwrap().to_string(), "false");
    }
}
