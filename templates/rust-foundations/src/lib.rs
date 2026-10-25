//! Reusable application features and their command-line adapter.
//!
//! Keep business rules in feature modules; the executable only connects I/O.
#![warn(missing_docs)]

pub mod cli;
pub mod features;

/// Public organisation/application identity and report metadata.
pub mod identity;
