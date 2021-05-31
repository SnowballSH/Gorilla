#[cfg(test)]
mod test_overall {
    use crate::builtin_types::integer::new_integer;
    use crate::builtin_types::string::new_string;
    use crate::helpers::run_code;

    #[test]
    fn connect() {
        let code = "1453 + 26 * 7";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(1453 + 26 * 7))));

        let code = "abc = 4763 / 23 + 98765; abc";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(4763 / 23 + 98765))));

        let code = "$abc = -9 * +-7; _ = -256 + $abc + 0; _";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(-256 + -9 * -7 + 0))));

        let code = "1 / 0";
        let result = run_code(code);
        assert!(result.is_err());

        let code = "1.this_does_not_exist";
        let result = run_code(code);
        assert!(result.is_err());

        let code = "a = true.i!.s!; a.i! - 1.i!";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(0))));

        let code = "\"01234\".i!";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(1234))));

        let code = "\"new\\nline\\0\"";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_string("new\nline\0".parse().unwrap()))));

        let code = "if 0
y
else
z";
        let result = run_code(code);
        assert!(result.is_err());

        let code = "result = 0
counter = 5
while counter > 0 {
    result = result + counter
    counter = counter - 1
}
result
";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(15))));

        let code = "fn add_anything(item1, item2) item1 + item2
add_anything(\"LOL \", \"ALR\")";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_string("LOL ALR".parse().unwrap()))));

        let code = "fn add_anything(item1, item2) item1 + item2
add_anything(75, 57)";
        let result = run_code(code);
        assert_eq!(result, Ok(Some(new_integer(75 + 57))));
    }
}