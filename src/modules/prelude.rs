use crate::builtin_types::null::new_null;
use crate::builtin_types::string::new_string;
use crate::env::Environment;
use crate::obj::{BaseObject, ObjResult};

pub fn print_line<'a>(
    _this: BaseObject<'a>,
    args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let mut strings = vec![];
    for arg in args {
        strings.push(arg.to_string());
    }
    let string = strings.join(" ");
    println!("{}", string);
    Ok(new_string(string))
}

pub fn puts<'a>(
    _this: BaseObject<'a>,
    args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    for arg in args {
        println!("{}", arg.to_string());
    }
    Ok(new_null())
}

pub fn print_inspect_line<'a>(
    _this: BaseObject<'a>,
    args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let mut strings = vec![];
    for arg in args {
        strings.push(arg.to_inspect_string());
    }
    let string = strings.join(" ");
    println!("{}", string);
    Ok(new_string(string))
}

pub fn print<'a>(
    _this: BaseObject<'a>,
    args: Vec<BaseObject<'a>>,
    _: Environment<'a>,
) -> ObjResult<'a> {
    let mut strings = vec![];
    for arg in args {
        strings.push(arg.to_string());
    }
    let string = strings.join(" ");
    print!("{}", string);
    Ok(new_string(string))
}
