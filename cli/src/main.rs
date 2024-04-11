use serde_aux::prelude::*;
use serde_json;
use std::error::Error;

use clap::Parser;
use reqwest::Url;

#[derive(serde::Serialize, serde::Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
struct Note {
    #[serde(deserialize_with = "deserialize_number_from_string")]
    id: i32,
    title: String,
    content: String,
    remind_date: String,
    created_at: String,
    updated_at: String,
}

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
    get_notes();
    println!("pattern: {:?}, path: {:?}", args.pattern, args.path)
}

fn get_notes() -> Result<(), Box<dyn Error>> {
    let url = Url::parse("http://localhost:8000/notes")?;
    let client = reqwest::blocking::Client::new();
    let resp = client
        .get(url)
        .query(&[("api-key", "123")])
        .send()
        .unwrap()
        .text()?;

    let notes: Vec<Note> = serde_json::from_str(&resp).unwrap();
    println!("{:?}", notes);
    Ok(())
}
