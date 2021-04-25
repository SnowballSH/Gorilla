use std::collections::HashMap;

use crate::obj::*;

#[derive(Default, Clone)]
pub(crate) struct Environment<'a> {
    pub(crate) store: HashMap<String, BaseObject<'a>>
}

impl<'a> Environment<'a> {
    pub(crate) fn set(&mut self, name: String, val: BaseObject<'a>) {
        self.store.insert(name, val);
    }

    pub(crate) fn get(&self, name: String) -> Option<BaseObject<'a>> {
        match self.store.get(&*name) {
            Some(x) => Some(x.clone()),
            None => None,
        }
    }
}
