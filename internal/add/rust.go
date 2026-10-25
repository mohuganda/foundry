package add

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RustFeature creates an isolated module; the caller explicitly registers its API.
func RustFeature(root, name string, options ...Options) (Result, error) {
	name, err := featureName(name)
	if err != nil {
		return Result{}, err
	}
	keywords := " as async await break const continue crate dyn else enum extern false fn for if impl in let loop match mod move mut pub ref return self static struct super trait true type unsafe use where while abstract become box do final gen macro override priv typeof unsized virtual yield try union "
	if strings.Contains(keywords, " "+name+" ") {
		return Result{}, fmt.Errorf("feature name %q is a Rust keyword; choose a name such as %s_feature", name, name)
	}
	if len(options) > 0 && options[0].App != "" {
		return Result{}, fmt.Errorf("--app is only supported for Flutter features")
	}
	base := filepath.Join("src", "features", name)
	files := map[string]string{
		filepath.Join(base, "mod.rs"): "//! " + featureTitle(name) + " business rules and application operations.\n\nmod domain;\nmod service;\n\npub use domain::{Error, Value};\npub use service::execute;\n",
		filepath.Join(base, "domain.rs"): `use std::error::Error as StdError;
use std::fmt;

/// An example validated feature value. Replace these rules with your domain rules.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Value(String);

impl Value {
    /// Trim surrounding whitespace and reject empty or control-containing input.
    ///
    /// # Errors
    /// Returns [Error] if the input violates the feature's validation rules.
    pub fn parse(input: &str) -> Result<Self, Error> {
        let input = input.trim();
        if input.is_empty() {
            return Err(Error::Empty);
        }
        if input.chars().any(char::is_control) {
            return Err(Error::ControlCharacter);
        }
        Ok(Self(input.to_owned()))
    }

    /// Borrow the validated value without exposing mutation.
    pub fn as_str(&self) -> &str {
        &self.0
    }
}

/// Reasons the example feature input cannot be accepted.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Error {
    /// The trimmed input is empty.
    Empty,
    /// The input contains a control character.
    ControlCharacter,
}

impl fmt::Display for Error {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        formatter.write_str(match self {
            Self::Empty => "value must not be empty",
            Self::ControlCharacter => "value must not contain control characters",
        })
    }
}

impl StdError for Error {}

#[cfg(test)]
mod tests {
    use super::{Error, Value};

    #[test]
    fn accepts_unicode_and_trims_surrounding_whitespace() {
        let value = Value::parse("  Zoë 李  ").expect("valid input");
        assert_eq!(value.as_str(), "Zoë 李");
    }

    #[test]
    fn rejects_empty_and_control_containing_input() {
        assert_eq!(Value::parse("   "), Err(Error::Empty));
        assert_eq!(Value::parse("one\ntwo"), Err(Error::ControlCharacter));
    }
}
`,
		filepath.Join(base, "service.rs"): `use super::{Error, Value};

/// Execute the feature's example operation using validated domain input.
///
/// Replace this body with product behavior; keep transport and I/O at boundaries.
///
/// # Errors
/// Returns [Error] when the input violates domain validation rules.
pub fn execute(input: &str) -> Result<Value, Error> {
    Value::parse(input)
}

#[cfg(test)]
mod tests {
    use super::{Error, execute};

    #[test]
    fn returns_a_validated_value() {
        let value = execute("  example  ").expect("valid input");
        assert_eq!(value.as_str(), "example");
    }

    #[test]
    fn propagates_validation_errors() {
        assert_eq!(execute(""), Err(Error::Empty));
    }
}
`,
		filepath.Join(base, "README.md"): "# " + featureTitle(name) + "\n\n" +
			"This module owns its domain validation and application operation. Replace the\n" +
			"example Value/Error types and execute function with your product's business rules.\n\n" +
			"Register it by adding `pub mod " + name + ";` to `src/features/mod.rs`. If the\n" +
			"application currently uses an inline features module, declare it there instead.\n" +
			"Then call `features::" + name + "::execute` from your application boundary.\n\n" +
			"The generator does not edit existing modules, CLI commands, Cargo.toml, or\n" +
			"Cargo.lock. The new unit tests run only after the module is registered.\n" +
			"Run cargo fmt --all and the project's full check script after wiring it.\n",
	}
	written, err := writeFiles(root, files, "rust", options...)
	if err != nil {
		return Result{}, err
	}
	return Result{Files: written, Next: []string{
		"register the module: add pub mod " + name + "; to src/features/mod.rs (or the existing inline features module)",
		"implement business rules and call features::" + name + "::execute from the application boundary",
		"run: cargo fmt --all, then bash scripts/check.sh (PowerShell: ./scripts/check.ps1)",
		"unit tests are compiled and run only after module registration",
	}}, nil
}
