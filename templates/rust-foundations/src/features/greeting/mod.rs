//! A small example feature with validated input and typed errors.

mod domain;
mod service;

pub use domain::{NameError, PersonName};
pub use service::greet;
