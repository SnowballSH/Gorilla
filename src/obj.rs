use std::any::Any;
use crate::env::*;

type CallFuncType =
fn(this: &BaseObject, args: Vec<&BaseObject>) ->
ObjResult;

pub type ObjResult = Result<&'static BaseObject, String>;
pub type ObjOption = Option<&'static BaseObject>;

#[derive(Clone, Copy)]
pub(crate) struct BaseObject {
    class: &'static Class,
    internal_value: &'static dyn Any,
    to_string_func: fn(this: &BaseObject) -> String,
    to_inspect_func: fn(this: &BaseObject) -> String,
    is_truthy_func: fn(this: &BaseObject) -> bool,
    equal_func: fn(this: &BaseObject, other: &BaseObject) -> bool,
    call_func: CallFuncType,
    parent_obj: &'static BaseObject,
}

impl BaseObject {
    #[inline]
    fn class(&self) -> &Class {
        self.class
    }

    fn value(&self) -> &dyn Any {
        self.internal_value
    }

    fn to_string(&self) -> String {
        (self.to_string_func)(self)
    }

    fn to_inspect_string(&self) -> String {
        (self.to_inspect_func)(self)
    }

    fn instance_get(&self, name: String) -> ObjOption {
        self.class.get_instance_var(name)
    }

    #[inline]
    fn is_truthy(&self) -> bool {
        (self.is_truthy_func)(self)
    }

    fn equal_to(&self, other: &BaseObject) -> bool {
        (self.equal_func)(self, other)
    }

    fn call(&self, this: &BaseObject, args: Vec<&BaseObject>) -> Result<&BaseObject, String> {
        (self.call_func)(this, args)
    }

    fn parent(&self) -> &BaseObject {
        self.parent_obj
    }

    fn set_parent(&mut self, parent: &'static BaseObject) {
        self.parent_obj = parent
    }
}

#[derive(Clone)]
struct Class {
    name: String,
    instance_vars: &'static Environment,
    super_class: &'static Class,
}

impl Class {
    fn to_string(&self) -> String {
        "Class '".to_owned() + &*self.name + "'"
    }

    fn to_inspect_string(&self) -> String {
        self.to_string()
    }

    fn parent(&self) -> &'static Class {
        self.super_class
    }

    fn set_parent(&mut self, _: &'static BaseObject) {}
    
    fn get_instance_var(&self, s: String) -> ObjOption {
        let x = self.instance_vars.get(s.clone());
        match x {
            Some(_) => x,
            None => self.super_class.get_instance_var(s)
        }
    }
}
