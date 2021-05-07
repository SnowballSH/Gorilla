use pest::Span;

#[inline]
pub fn span_to_line(source: &str, span: Span) -> usize {
    source[..span.start()].matches("\n").count()
}

pub fn leb128_unsigned(val: u64) -> Vec<u8> {
    let mut value = val | 0;
    if value < 0x80 {
        return vec![value as u8];
    }

    let mut res = vec![];

    loop {
        let mut c = (value & 0x7f) as u8;
        value >>= 7;
        if value != 0 {
            c |= 0x80;
        }
        res.push(c);
        if c & 0x80 == 0 {
            break;
        }
    }

    res
}
