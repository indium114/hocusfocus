use dirs;

// MARK: dir helpers
pub fn home_dir() -> String {
    let dir = dirs::home_dir();
    return dir
        .map(|p| p.to_string_lossy().into_owned())
        .unwrap_or_default();
}

// MARK: subcommand helpers
pub fn print_help() {
    println!("hocusfocus help:");
    println!(" <no args>      : choose/stop session");
    println!(" help           : print this message");
    println!(" currentsession : print current session");
    println!(" stats          : print statistics");
}
