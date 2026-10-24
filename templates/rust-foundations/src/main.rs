use std::env;
use std::io;
use std::process::ExitCode;

fn main() -> ExitCode {
    match __RUST_CRATE__::cli::run(env::args().skip(1), &mut io::stdout().lock()) {
        Ok(()) => ExitCode::SUCCESS,
        Err(error) => {
            eprintln!("error: {error}");
            ExitCode::FAILURE
        }
    }
}
