#![forbid(unsafe_code)]

use crate::env::*;
use crate::obj::ValueType::*;

pub type CallFuncType<'a> = fn(BaseObject<'a>, Vec<BaseObject<'a>>, Environment<'a>) -> ObjResult<'a>;

#[inline]
pub fn not_callable<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>, _: Environment) -> ObjResult<'a> {
    Err(format!(
        "'{}' ({}) is not callable",
        this.to_string(),
        this.class.to_string()
    ))
}

pub type ObjResult<'a> = Result<BaseObject<'a>, String>;

pub type NativeFunctionType<'a> = (&'static str, CallFuncType<'a>);
pub type FunctionType = (String, Vec<String>, Vec<u8>);

#[derive(Clone, Eq, Debug)]
pub enum ValueType<'a> {
    Int(i64),
    Bool(bool),
    NativeFunction(NativeFunctionType<'a>),
    Function(FunctionType),
    Str(String),
    Vector(Vec<&'a BaseObject<'a>>),
    Null,
}

impl<'a> PartialEq for ValueType<'a> {
    fn eq(&self, other: &Self) -> bool {
        match self {
            Int(a) => match other {
                Int(i) => a == i,
                _ => false,
            },
            Bool(a) => match other {
                Bool(i) => a == i,
                _ => false,
            },
            Str(a) => match other {
                Str(i) => a == i,
                _ => false,
            },
            NativeFunction(a) => match other {
                NativeFunction(i) => a == i,
                _ => false,
            },
            Function(a) => match other {
                Function(i) => a == i,
                _ => false,
            },
            Vector(a) => match other {
                Vector(i) => a == i,
                _ => false,
            },
            Null => match other {
                Null => true,
                _ => false,
            },
        }
    }
}

#[derive(Clone, Eq, PartialEq, Debug)]
pub struct BaseObject<'a> {
    pub class: Class<'a>,
    pub internal_value: ValueType<'a>,
    pub to_string_func: fn(BaseObject<'a>) -> String,
    pub to_inspect_func: fn(BaseObject<'a>) -> String,
    pub is_truthy_func: fn(BaseObject<'a>) -> bool,
    pub equal_func: fn(BaseObject<'a>, BaseObject<'a>) -> bool,
    pub call_func: CallFuncType<'a>,
    pub parent_obj: Option<Box<BaseObject<'a>>>,
}

impl<'a> BaseObject<'a> {
    #[inline]
    pub fn to_string(&self) -> String {
        (self.to_string_func)(self.clone())
    }

    #[inline]
    pub fn to_inspect_string(&self) -> String {
        (self.to_inspect_func)(self.clone())
    }

    #[inline]
    pub fn instance_get(&self, name: String) -> Option<BaseObject<'a>> {
        self.class.get_instance_var(name)
    }

    #[inline]
    pub fn is_truthy(&self) -> bool {
        (self.is_truthy_func)(self.clone())
    }

    #[inline]
    pub fn equal_to(&self, other: BaseObject<'a>) -> bool {
        (self.equal_func)(self.clone(), other)
    }

    pub fn call(
        &self,
        this: BaseObject<'a>,
        args: Vec<BaseObject<'a>>,
        vm: Environment<'a>,
    ) -> Result<BaseObject<'a>, String> {
        (self.call_func)(this, args, vm)
    }

    #[inline]
    pub fn parent(&self) -> Option<Box<BaseObject<'a>>> {
        self.clone().parent_obj
    }

    pub fn set_parent(&mut self, parent: BaseObject<'a>) {
        self.parent_obj = Some(Box::from(parent))
    }
}

#[derive(Clone, Eq, Debug)]
pub struct Class<'a> {
    pub name: &'a str,
    pub instance_vars: Environment<'a>,
    pub super_class: Option<Box<Class<'a>>>,
}

impl<'a> Class<'a> {
    #[inline]
    pub fn to_string(&self) -> String {
        "Class '".to_owned() + self.name + "'"
    }

    pub fn get_instance_var(&self, s: String) -> Option<BaseObject<'a>> {
        let x = self.instance_vars.get(s.clone());
        match x {
            Some(_) => x,
            None => match self.super_class.clone() {
                Some(y) => y.get_instance_var(s),
                None => None,
            },
        }
    }
}

impl<'a> PartialEq for Class<'a> {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name && self.super_class == other.super_class
    }
}
