use super::{NameError, PersonName};

/// Build a greeting after validating its recipient.
///
/// ```
/// use __RUST_CRATE__::features::greeting::greet;
///
/// assert_eq!(greet("Ada").expect("valid name"), "Hello, Ada!");
/// ```
///
/// # Errors
/// Returns [`NameError`] when the recipient is not a valid name.
pub fn greet(recipient: &str) -> Result<String, NameError> {
    let name = PersonName::parse(recipient)?;
    Ok(format!("Hello, {}!", name.as_str()))
}
