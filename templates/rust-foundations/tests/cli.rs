use std::process::Command;

fn command() -> Command {
    Command::new(env!("CARGO_BIN_EXE___RUST_CRATE__"))
}

#[test]
fn greets_a_recipient_from_the_executable() {
    let output = command()
        .args(["greet", "Ada"])
        .output()
        .expect("run binary");
    assert!(output.status.success());
    assert_eq!(
        String::from_utf8(output.stdout).expect("UTF-8 output"),
        "Hello, Ada!\n"
    );
    assert!(output.stderr.is_empty());
}

#[test]
fn displays_help_successfully() {
    let output = command().arg("--help").output().expect("run binary");
    assert!(output.status.success());
    assert!(
        String::from_utf8(output.stdout)
            .expect("UTF-8 output")
            .contains("greet <name>")
    );
}

#[test]
fn rejects_bad_commands_and_invalid_names() {
    for args in [
        vec!["unknown"],
        vec!["greet"],
        vec!["greet", "  "],
        vec!["greet", "Ada", "extra"],
    ] {
        let output = command().args(args).output().expect("run binary");
        assert!(!output.status.success());
        assert!(output.stdout.is_empty());
        assert!(
            String::from_utf8(output.stderr)
                .expect("UTF-8 error")
                .contains("error:")
        );
    }
}
