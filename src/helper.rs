pub fn wrap<T>(x: T) -> Box<T> {
    Box::new(x)
}