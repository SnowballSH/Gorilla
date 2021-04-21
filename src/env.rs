use std::collections::HashMap;
use crate::obj::*;

type StoreType = HashMap<String, BaseObject>;

pub struct Environment {
    store: StoreType
}

impl Environment {
    pub fn set(&'static mut self, name: String, val: BaseObject) {
        self.store.insert(name, val);
    }

    pub fn get(&'static self, name: String) -> ObjOption {
        self.store.get(&*name)
    }
}
