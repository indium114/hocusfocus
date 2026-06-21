use demand::{DemandOption, Select};
use std::env;

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
                println!("called currentsession");
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

    let sessions: Vec<help::Session> = help::load_sessions();

    let form = Select::new("HocusFocus")
        .description("choose a session type")
        .filterable(true)
        .option(DemandOption::new("Work"))
        .option(DemandOption::new("Study"))
        .option(DemandOption::new("Waste"))
        .option(DemandOption::new("Stop Current Session"));
    let _ = form.run();

    println!("pretend this is a menu")
}
