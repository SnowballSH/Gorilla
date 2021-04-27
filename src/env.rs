use std::collections::HashMap;

use crate::obj::*;

#[derive(Default, Clone, Eq, PartialEq, Debug)]
pub struct Environment<'a> {
    pub store: HashMap<String, BaseObject<'a>>
}

impl<'a> Environment<'a> {
    pub fn set(&mut self, name: String, val: BaseObject<'a>) {
        self.store.insert(name, val);
    }

    pub fn get(&self, name: String) -> Option<BaseObject<'a>> {
        match self.store.get(&*name) {
            Some(x) => Some(x.clone()),
            None => None,
        }
    }
}
