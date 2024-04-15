use serde_aux::prelude::*;
use serde_json;
use std::error::Error;

use reqwest::Url;

#[derive(serde::Serialize, serde::Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct Note {
    #[serde(deserialize_with = "deserialize_number_from_string")]
    id: i32,
    title: String,
    content: String,
    remind_date: String,
    created_at: String,
    updated_at: String,
}

#[derive(serde::Serialize, serde::Deserialize, Debug)]
#[serde(rename_all = "camelCase")]
pub struct NoteCreate {
    title: String,
    content: String,
    remind_date: String,
}

pub fn get_notes() -> Result<Vec<Note>, Box<dyn Error>> {
    let url = Url::parse("http://localhost:8000/notes")?;
    let client = reqwest::blocking::Client::new();
    let resp = client
        .get(url)
        .query(&[("api-key", "123")])
        .send()
        .unwrap()
        .text()?;

    let notes: Vec<Note> = serde_json::from_str(&resp).unwrap();
    Ok(notes)
}

pub fn create_note(title: &str, content: &str, remind_date: &str) -> Result<Note, Box<dyn Error>> {
    let url = Url::parse("http://localhost:8000/notes")?;
    let client = reqwest::blocking::Client::new();
    let note = NoteCreate {
        title: title.to_string(),
        content: content.to_string(),
        remind_date: remind_date.to_string(),
    };

    let resp = client
        .post(url)
        .query(&[("api-key", "123")])
        .header("Content-Type", "application/json")
        .body(serde_json::to_string(&note)?)
        .send()
        .unwrap()
        .text()?;

    let note: Note = serde_json::from_str(&resp).unwrap();
    Ok(note)
}

pub fn update_note(
    id: i32,
    title: &str,
    content: &str,
    remind_date: &str,
) -> Result<Note, Box<dyn Error>> {
    let url = Url::parse(&format!("http://localhost:8000/notes/{}", id))?;
    let client = reqwest::blocking::Client::new();
    let note = NoteCreate {
        title: title.to_string(),
        content: content.to_string(),
        remind_date: remind_date.to_string(),
    };

    let resp = client
        .put(url)
        .query(&[("api-key", "123")])
        .header("Content-Type", "application/json")
        .body(serde_json::to_string(&note)?)
        .send()
        .unwrap()
        .text()?;

    let note: Note = serde_json::from_str(&resp).unwrap();
    Ok(note)
}

pub fn delete_note(id: i32) -> Result<(), Box<dyn Error>> {
    let url = Url::parse(&format!("http://localhost:8000/notes/{}", id))?;
    let client = reqwest::blocking::Client::new();
    client
        .delete(url)
        .query(&[("api-key", "123")])
        .send()
        .unwrap();
    Ok(())
}
