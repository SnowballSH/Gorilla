param([String]$p1);
go build;
./Gorilla -o a.grx -c $p1;
cargo run --quiet -- a.grx;
rm a.grx;