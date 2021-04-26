use crate::env::*;

pub type CallFuncType<'a> =
fn(BaseObject<'a>, Vec<BaseObject<'a>>) -> ObjResult<'a>;

#[inline]
pub fn not_callable<'a>() -> CallFuncType<'a> {
    fn a<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
        Err(format!("'{}' ({}) is not callable", this.to_string(), this.class.to_string()))
    }

    a
}

pub type ObjResult<'a> = Result<BaseObject<'a>, String>;

#[derive(Copy)]
pub union ValueType {
    pub int: i64,
}

impl Clone for ValueType {
    #[inline]
    fn clone(&self) -> Self {
        unsafe {
            match self {
                ValueType { int } => ValueType { int: int.clone() }
            }
        }
    }
}

impl PartialEq for ValueType {
    fn eq(&self, other: &Self) -> bool {
        unsafe {
            match self {
                ValueType { int } => {
                    let a = int;
                    match other {
                        ValueType { int } => a == int,
                    }
                }
                _ => false
            }
        }
    }
}

#[derive(Clone)]
pub struct BaseObject<'a> {
    pub class: Class<'a>,
    pub internal_value: ValueType,
    pub to_string_func: fn(BaseObject<'a>) -> String,
    pub to_inspect_func: fn(BaseObject<'a>) -> String,
    pub is_truthy_func: fn(BaseObject<'a>) -> bool,
    pub equal_func: fn(BaseObject<'a>, BaseObject<'a>) -> bool,
    pub call_func: CallFuncType<'a>,
    pub parent_obj: Option<Box<BaseObject<'a>>>,
}

impl<'a> BaseObject<'a> {
    #[inline]
    pub fn class(&self) -> &Class<'a> {
        &self.class
    }

    #[inline]
    pub fn value(&self) -> ValueType {
        self.internal_value
    }

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

    pub fn call(&self, this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> Result<BaseObject<'a>, String> {
        (self.call_func)(this, args)
    }

    #[inline]
    pub fn parent(&self) -> Option<Box<BaseObject<'a>>> {
        self.clone().parent_obj
    }

    pub fn set_parent(&mut self, parent: BaseObject<'a>) {
        self.parent_obj = Some(Box::from(parent))
    }
}

#[derive(Clone)]
pub struct Class<'a> {
    pub name: &'static str,
    pub instance_vars: Environment<'a>,
    pub super_class: Option<&'a Class<'a>>,
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
            None => match self.super_class {
                Some(y) => y.get_instance_var(s),
                None => None
            }
        }
    }
}

impl<'a> PartialEq for Class<'a> {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name && self.super_class == other.super_class
    }
}
