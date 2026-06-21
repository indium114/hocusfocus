use std::env;

mod help;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() > 1 {
        match args[1] {
            "help" => {
                println!("called help")
            }
            "currentsession" => {
                println!("called currentsession")
            }
            "stats" => {
                println!("called stats")
            }
            _ => {
                println!("called help through unknown arg")
            }
        }
    }

    println!("pretend this is a menu")
}
