use std::env;

mod help;

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() > 1 {
        match args[1].as_str() {
            "help" => {
                println!("called help")
                return
            }
            "currentsession" => {
                println!("called currentsession")
                return
            }
            "stats" => {
                println!("called stats")
                return
            }
            _ => {
                println!("called help through unknown arg")
                return
            }
        }
    }

    println!("pretend this is a menu")
}
