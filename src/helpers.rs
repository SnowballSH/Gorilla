use pest::Span;

pub fn span_to_line(source: &str, span: Span) -> usize {
    source[..span.start()].matches("\n").count()
}