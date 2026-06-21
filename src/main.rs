use chrono::Local;
use demand::{DemandOption, Select};
use humantime::format_duration;
use std::env;
use std::time::Duration;

mod help;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() > 1 {
        match args[1].as_str() {
            "help" => {
                help::print_help();
                return;
            }
            "currentsession" => {
                let sessions: Vec<help::Session> = help::load_sessions();
                let current: Option<&help::Session> = help::current_session(&sessions);

                match current {
                    Some(session) => {
                        let elapsed = Local::now().signed_duration_since(session.start);
                        let elapsed_dur = Duration::from_secs(elapsed.num_seconds() as u64);
                        let format_elapsed = format_duration(elapsed_dur);
                        println!(" Current session: {} ({})", session.kind, format_elapsed);
                        return;
                    }
                    None => {
                        println!(" No current session")
                    }
                }

                return;
            }
            "stats" => {
                println!("called stats");
                return;
            }
            _ => {
                println!("called help through unknown arg");
                return;
            }
        }
    }

    let mut sessions: Vec<help::Session> = help::load_sessions();

    let form = Select::new("HocusFocus")
        .description("choose a session type")
        .filterable(true)
        .option(DemandOption::new("Work"))
        .option(DemandOption::new("Study"))
        .option(DemandOption::new("Waste"))
        .option(DemandOption::new("Stop Current Session"));

    match form.run() {
        Ok(choice) => match choice {
            "Work" => {
                let _ = help::stop_session(&mut sessions);
                help::start_session("Work".to_string(), &mut sessions);
                help::save_sessions(sessions);
            }
            "Study" => {
                let _ = help::stop_session(&mut sessions);
                help::start_session("Study".to_string(), &mut sessions);
                help::save_sessions(sessions);
            }
            "Waste" => {
                let _ = help::stop_session(&mut sessions);
                help::start_session("Waste".to_string(), &mut sessions);
                help::save_sessions(sessions);
            }
            "Stop Current Session" => {
                let current = help::current_session(&sessions);
                match current {
                    Some(_) => {
                        help::stop_session(&mut sessions);
                        let success = help::save_sessions(sessions);
                        match success {
                            true => (),
                            false => {
                                println!("failed to save sessions");
                            }
                        }
                    }
                    None => {
                        println!("No current session")
                    }
                }
            }
            _ => (),
        },
        Err(e) => {
            panic!("error: {}", e);
        }
    }
}
