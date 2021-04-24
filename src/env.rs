use std::collections::HashMap;

use crate::obj::*;

type StoreType = HashMap<String, BaseObject>;

#[derive(Default)]
pub(crate) struct Environment {
    pub(crate) store: StoreType
}

impl Environment {
    pub(crate) fn set(&'static mut self, name: String, val: BaseObject) {
        self.store.insert(name, val);
    }

    pub(crate) fn get(&'static self, name: String) -> ObjOption {
        self.store.get(&*name)
    }
}
