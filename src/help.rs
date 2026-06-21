use chrono::{DateTime, FixedOffset, Local};
use dirs;
use humantime::format_duration;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::time::{Duration, Instant};

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

pub fn stop_session(sessions: &mut [Session]) {
    let now = Local::now();

    for idx in sessions {
        if idx.end.is_none() {
            let current: DateTime<FixedOffset> = now.into();
            idx.end = Some(current);
        }
    }
}

pub fn start_session(kind: String, sessions: &mut Vec<Session>) {
    let now: DateTime<FixedOffset> = Local::now().into();

    let new_session = Session {
        kind: kind,
        start: now,
        end: None,
    };

    sessions.push(new_session);
}

pub fn print_stats() {
    let sessions = load_sessions();
    let mut totals: HashMap<String, Duration> = HashMap::new();

    for i in sessions {
        if let Some(end_time) = i.end {
            let dur = end_time - i.start;
            let seconds_dur = Duration::from_secs(dur.num_seconds() as u64);
            let formatted_dur = format_duration(seconds_dur).to_string();
            let std_dur = formatted_dur.parse::<humantime::Duration>().unwrap().into();
            *totals.entry(i.kind).or_insert(Duration::ZERO) += std_dur;
        } else {
            continue;
        }
    }

    if totals.len() == 0 {
        println!(" No sessions have been completed.");
        return;
    }

    for (kind, dur) in totals {
        println!("{}: {}", kind, format_duration(dur))
    }
}
