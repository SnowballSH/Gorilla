#[cfg(test)]
mod test_overall {
    use crate::integer::new_integer;
    use crate::helpers::run_code;

    #[test]
    fn connect() {
        let code = "1453 + 26 * 7";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(1453 + 26 * 7))));

        let code = "abc = 4763 / 23 + 98765; abc";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(4763 / 23 + 98765))));
    }
}