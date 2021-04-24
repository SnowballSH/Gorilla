use crate::env::*;

pub(crate) type CallFuncType =
fn(&BaseObject, Vec<&BaseObject>) -> ObjResult;

#[inline]
pub(crate) fn not_callable() -> CallFuncType {
    fn a (this: &BaseObject, _args: Vec<&BaseObject>) -> ObjResult {
        Err(format!("'{}' ({}) is not callable", this.to_string(), this.class.to_string()))
    }

    a
}

pub(crate) type ObjResult = Result<&'static BaseObject, String>;
pub(crate) type ObjOption = Option<&'static BaseObject>;

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
pub(crate) struct BaseObject {
    pub(crate) class: Class,
    pub(crate) internal_value: ValueType,
    pub(crate) to_string_func: fn(&BaseObject) -> String,
    pub(crate) to_inspect_func: fn(&BaseObject) -> String,
    pub(crate) is_truthy_func: fn(&BaseObject) -> bool,
    pub(crate) equal_func: fn(&BaseObject, &BaseObject) -> bool,
    pub(crate) call_func: CallFuncType,
    pub(crate) parent_obj: ObjOption,
}

impl BaseObject {
    #[inline]
    fn class(&self) -> &Class {
        &self.class
    }

    #[inline]
    fn value(&self) -> ValueType {
        self.internal_value
    }

    #[inline]
    fn to_string(&self) -> String {
        (self.to_string_func)(self)
    }

    #[inline]
    fn to_inspect_string(&self) -> String {
        (self.to_inspect_func)(self)
    }

    #[inline]
    fn instance_get(& self, name: String) -> Option<BaseObject> {
        self.class.get_instance_var(name)
    }

    #[inline]
    fn is_truthy(&self) -> bool {
        (self.is_truthy_func)(self)
    }

    #[inline]
    fn equal_to(&self, other: &BaseObject) -> bool {
        (self.equal_func)(self, other)
    }

    fn call(&self, this: &BaseObject, args: Vec<&BaseObject>) -> Result<&BaseObject, String> {
        (self.call_func)(this, args)
    }

    #[inline]
    fn parent(&self) -> ObjOption {
        self.parent_obj
    }

    fn set_parent(&mut self, parent: &'static BaseObject) {
        self.parent_obj = Some(parent)
    }
}

#[derive(Clone)]
pub(crate) struct Class {
    pub(crate) name: &'static str,
    pub(crate) instance_vars: Environment,
    pub(crate) super_class: Option<&'static Class>,
}

impl Class {
    #[inline]
    fn to_string(&self) -> String {
        "Class '".to_owned() + self.name + "'"
    }

    #[inline]
    fn to_inspect_string(&self) -> String {
        self.to_string()
    }

    fn set_parent(&mut self, _: &'static BaseObject) {}

    fn get_instance_var(&self, s: String) -> Option<BaseObject> {
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

impl PartialEq for Class {
    fn eq(&self, other: &Self) -> bool {
        self.name == other.name && self.super_class == other.super_class
    }
}
