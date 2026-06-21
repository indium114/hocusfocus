use dirs;

fn home_dir() -> String {
    let dir = dirs::home_dir();
    return dir
        .map(|p| p.to_string_lossy().into_owned())
        .unwrap_or_default();
}
