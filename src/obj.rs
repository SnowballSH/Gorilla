use crate::env::*;

pub(crate) type CallFuncType<'a> =
fn(BaseObject<'a>, Vec<BaseObject<'a>>) -> ObjResult<'a>;

#[inline]
pub(crate) fn not_callable<'a>() -> CallFuncType<'a> {
    fn a<'a>(this: BaseObject<'a>, _args: Vec<BaseObject<'a>>) -> ObjResult<'a> {
        Err(format!("'{}' ({}) is not callable", this.to_string(), this.class.to_string()))
    }

    a
}

pub(crate) type ObjResult<'a> = Result<BaseObject<'a>, String>;
pub(crate) type ObjOption<'a> = Option<BaseObject<'a>>;

#[derive(Copy)]
pub(crate) union ValueType {
    pub(crate) int: i64,
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
                        _ => false
                    }
                }
                _ => false
            }
        }
    }
}

#[derive(Clone)]
pub(crate) struct BaseObject<'a> {
    pub(crate) class: Class<'a>,
    pub(crate) internal_value: ValueType,
    pub(crate) to_string_func: fn(BaseObject<'a>) -> String,
    pub(crate) to_inspect_func: fn(BaseObject<'a>) -> String,
    pub(crate) is_truthy_func: fn(BaseObject<'a>) -> bool,
    pub(crate) equal_func: fn(BaseObject<'a>, BaseObject<'a>) -> bool,
    pub(crate) call_func: CallFuncType<'a>,
    pub(crate) parent_obj: Option<Box<BaseObject<'a>>>,
}

impl<'a> BaseObject<'a> {
    #[inline]
    pub(crate) fn class(&self) -> &Class<'a> {
        &self.class
    }

    #[inline]
    pub(crate) fn value(&self) -> ValueType {
        self.internal_value
    }

    #[inline]
    pub(crate) fn to_string(&self) -> String {
        (self.to_string_func)(self.clone())
    }

    #[inline]
    pub(crate) fn to_inspect_string(&self) -> String {
        (self.to_inspect_func)(self.clone())
    }

    #[inline]
    pub(crate) fn instance_get(&self, name: String) -> Option<BaseObject<'a>> {
        self.class.get_instance_var(name)
    }

    #[inline]
    pub(crate) fn is_truthy(&self) -> bool {
        (self.is_truthy_func)(self.clone())
    }

    #[inline]
    pub(crate) fn equal_to(&self, other: BaseObject<'a>) -> bool {
        (self.equal_func)(self.clone(), other)
    }

    pub(crate) fn call(&self, this: BaseObject<'a>, args: Vec<BaseObject<'a>>) -> Result<BaseObject<'a>, String> {
        (self.call_func)(this, args)
    }

    #[inline]
    pub(crate) fn parent(&self) -> Option<Box<BaseObject<'a>>> {
        self.clone().parent_obj
    }

    pub(crate) fn set_parent(&mut self, parent: BaseObject<'a>) {
        self.parent_obj = Some(Box::from(parent))
    }
}

#[derive(Clone)]
pub(crate) struct Class<'a> {
    pub(crate) name: &'static str,
    pub(crate) instance_vars: Environment<'a>,
    pub(crate) super_class: Option<&'a Class<'a>>,
}

impl<'a> Class<'a> {
    #[inline]
    pub(crate) fn to_string(&self) -> String {
        "Class '".to_owned() + self.name + "'"
    }

    pub(crate) fn get_instance_var(&self, s: String) -> Option<BaseObject<'a>> {
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
