use std::error::Error;
use std::fmt;

/// A trimmed name with no control characters and at most 80 Unicode characters.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct PersonName(String);

impl PersonName {
    /// Trim surrounding whitespace and validate a name.
    ///
    /// # Errors
    /// Returns [`NameError`] for empty, excessively long, or control-containing names.
    pub fn parse(value: &str) -> Result<Self, NameError> {
        let value = value.trim();
        if value.is_empty() {
            return Err(NameError::Empty);
        }
        if value.chars().count() > 80 {
            return Err(NameError::TooLong);
        }
        if value.chars().any(char::is_control) {
            return Err(NameError::ControlCharacter);
        }
        Ok(Self(value.to_owned()))
    }

    /// Borrow the validated name without exposing mutation.
    pub fn as_str(&self) -> &str {
        &self.0
    }
}

/// Reasons why a name cannot be used in a greeting.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum NameError {
    /// The trimmed name is empty.
    Empty,
    /// The name contains more than 80 Unicode characters.
    TooLong,
    /// The name includes a control character.
    ControlCharacter,
}

impl fmt::Display for NameError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        let message = match self {
            Self::Empty => "name must not be empty",
            Self::TooLong => "name must be at most 80 characters",
            Self::ControlCharacter => "name must not contain control characters",
        };
        formatter.write_str(message)
    }
}

impl Error for NameError {}

#[cfg(test)]
mod tests {
    use super::{NameError, PersonName};

    #[test]
    fn accepts_unicode_and_normalizes_surrounding_whitespace() {
        let name = PersonName::parse("  Zoë 李  ").expect("valid name");
        assert_eq!(name.as_str(), "Zoë 李");
        assert!(PersonName::parse(&"界".repeat(80)).is_ok());
    }

    #[test]
    fn rejects_empty_long_and_control_containing_names() {
        for (input, expected) in [
            ("  ".to_owned(), NameError::Empty),
            ("界".repeat(81), NameError::TooLong),
            ("Ada\nLovelace".to_owned(), NameError::ControlCharacter),
        ] {
            assert_eq!(PersonName::parse(&input), Err(expected));
        }
    }
}
