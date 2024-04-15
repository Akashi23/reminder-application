mod notes;

use clap::Parser;

/// Search for a pattern in a file and display the lines that contain it.
#[derive(Parser)]
struct Cli {
    /// The pattern to look for
    pattern: String,
    /// The path to the file to read
    path: std::path::PathBuf,
}

fn main() {
    let args = Cli::parse();
    // let notes = notes::update_note(10, "title 10", "Con", "2024/16/04");
    // let notes = notes::create_note("ti", "contentdsa1212", "2024/16/04");
    let notes = notes::get_notes().unwrap();
    // let notes = notes::delete_note(5);
    println!("pattern: {:?}, path: {:?}", args.pattern, args.path);
    println!("{:?}", notes);
}
