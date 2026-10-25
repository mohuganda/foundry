//! Command-line parsing and output, kept separate from business rules.

use std::error::Error;
use std::fmt;
use std::io::{self, Write};

use crate::features::greeting::{NameError, greet};

/// Failures at the command-line boundary.
#[derive(Debug)]
pub enum CliError {
    /// The supplied arguments do not match a supported command.
    Usage,
    /// A feature rejected its input.
    InvalidName(NameError),
    /// Writing command output failed.
    Io(io::Error),
}

impl fmt::Display for CliError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Usage => write!(formatter, "expected 'greet <name>'; use --help for usage"),
            Self::InvalidName(error) => error.fmt(formatter),
            Self::Io(error) => error.fmt(formatter),
        }
    }
}

impl Error for CliError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match self {
            Self::Usage => None,
            Self::InvalidName(error) => Some(error),
            Self::Io(error) => Some(error),
        }
    }
}

/// Run a command using explicit arguments and an output writer.
///
/// # Errors
/// Returns [`CliError`] for unsupported arguments, invalid names, or output errors.
pub fn run(
    args: impl IntoIterator<Item = String>,
    output: &mut impl Write,
) -> Result<(), CliError> {
    let args: Vec<String> = args.into_iter().collect();
    let args: Vec<&str> = args.iter().map(String::as_str).collect();
    let message = match args.as_slice() {
        [] | ["--help"] | ["-h"] | ["greet", "--help"] => {
            format!(
                "Usage: {} greet <name>\n       {} --help\n       identity (public application metadata)",
                env!("CARGO_PKG_NAME"),
                env!("CARGO_PKG_NAME")
            )
        }
        ["identity"] => crate::identity::PUBLIC_JSON.to_owned(),
        ["greet", name] if !name.starts_with('-') => greet(name).map_err(CliError::InvalidName)?,
        _ => return Err(CliError::Usage),
    };
    writeln!(output, "{message}").map_err(CliError::Io)
}

#[cfg(test)]
mod tests {
    use super::{CliError, run};
    use std::io::{self, Write};

    #[test]
    fn accepts_a_quoted_multi_word_name() {
        let mut output = Vec::new();
        run(["greet".to_owned(), "Ada Lovelace".to_owned()], &mut output).expect("valid command");
        assert_eq!(output, b"Hello, Ada Lovelace!\n");
    }

    #[test]
    fn returns_writer_errors_without_panicking() {
        struct BrokenWriter;
        impl Write for BrokenWriter {
            fn write(&mut self, _: &[u8]) -> io::Result<usize> {
                Err(io::ErrorKind::BrokenPipe.into())
            }
            fn flush(&mut self) -> io::Result<()> {
                Ok(())
            }
        }
        let error = run(["--help".to_owned()], &mut BrokenWriter).expect_err("writer fails");
        assert!(matches!(error, CliError::Io(error) if error.kind() == io::ErrorKind::BrokenPipe));
    }
}
