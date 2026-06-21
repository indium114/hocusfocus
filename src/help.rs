use chrono::{DateTime, FixedOffset, Local};
use dirs;
use serde::{Deserialize, Serialize};
use std::fs;

// MARK: types
#[derive(Serialize, Deserialize, Debug)]
pub struct Session {
    #[serde(rename = "type")]
    pub kind: String,
    pub start: DateTime<FixedOffset>,
    pub end: Option<DateTime<FixedOffset>>,
}

// MARK: dir helpers
pub fn home_dir() -> String {
    let dir = dirs::home_dir();
    return dir
        .map(|p| p.to_string_lossy().into_owned())
        .unwrap_or_default();
}

pub fn store_path() -> String {
    return home_dir() + "/.hocusfocus.json";
}

// MARK: save/load helpers
pub fn load_sessions() -> Vec<Session> {
    fs::read_to_string(store_path())
        .ok()
        .and_then(|s| serde_json::from_str(&s).ok())
        .unwrap_or_default()
}

pub fn save_sessions(sessions: Vec<Session>) -> bool {
    match serde_json::to_string_pretty(&sessions) {
        Ok(json) => fs::write(store_path(), json).is_ok(),
        Err(_) => false,
    }
}

// MARK: subcommand helpers
pub fn print_help() {
    println!("hocusfocus help:");
    println!(" <no args>      : choose/stop session");
    println!(" help           : print this message");
    println!(" currentsession : print current session");
    println!(" stats          : print statistics");
}

pub fn current_session(sessions: &[Session]) -> Option<&Session> {
    for session in sessions {
        if session.end.is_none() {
            return Some(session);
        }
    }
    return None;
}
