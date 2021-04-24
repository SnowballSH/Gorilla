use std::collections::HashMap;

use crate::obj::*;

type StoreType = HashMap<String, BaseObject>;

#[derive(Default, Clone)]
pub(crate) struct Environment {
    pub(crate) store: StoreType
}

impl Environment {
    pub(crate) fn set(&mut self, name: String, val: BaseObject) {
        self.store.insert(name, val);
    }

    pub(crate) fn get(&self, name: String) -> Option<BaseObject> {
        match self.store.get(&*name) {
            Some(x) => Some(x.to_owned()),
            None => None,
        }
    }
}
