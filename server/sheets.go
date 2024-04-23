package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Note struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	RemindDate string `json:"remindDate"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

var (
	spreadsheetID = GetEnv("SPREADSHEET_ID")
	readRange     = "Sheet1!A:F"
	credentials   = GetEnv("CREDENTIALS")
)

func getClient() *sheets.Service {
	creds, err := os.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(context.Background())

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv
}

func GetNote(id string) Note {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {
		if row[0] == id {
			return Note{
				Id:         row[0].(string),
				Title:      row[1].(string),
				Content:    row[2].(string),
				RemindDate: row[3].(string),
				CreatedAt:  row[4].(string),
				UpdatedAt:  row[5].(string),
			}
		}
	}

	return Note{}
}

func CreateNote(note Note) Note {
	sheetsService := getClient()

	note.Id = fmt.Sprint(getLastId() + 1)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{
			{note.Id, note.Title, note.Content, note.RemindDate, note.CreatedAt, note.UpdatedAt},
		},
	}

	_, err := sheetsService.Spreadsheets.Values.Append(spreadsheetID, readRange, valueRange).ValueInputOption("RAW").Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	return note
}

func UpdateNote(note Note) Note {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for i, row := range resp.Values {
		if row[0] == note.Id {
			createdAt := resp.Values[i][4].(string)
			resp.Values[i] = []interface{}{note.Id, note.Title, note.Content, note.RemindDate, createdAt, note.UpdatedAt}
			note.CreatedAt = createdAt
		}
	}

	valueRange := &sheets.ValueRange{
		Values: resp.Values,
	}

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID, readRange, valueRange).ValueInputOption("RAW").Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	return note
}

func DeleteNote(id string) Note {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for i, row := range resp.Values {
		if row[0] == id {
			resp.Values = append(resp.Values[:i], resp.Values[i+1:]...)
		}
	}

	valueRange := &sheets.ValueRange{
		Values: resp.Values,
	}

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID, readRange, valueRange).ValueInputOption("RAW").Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	resp, err = sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for i, row := range resp.Values {
		row[0] = fmt.Sprint(i + 1)
	}

	lastId := getLastId() - 1

	resp.Values[lastId] = []interface{}{"", "", "", "", "", ""}

	valueRange = &sheets.ValueRange{
		Values: resp.Values,
	}

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID, readRange, valueRange).ValueInputOption("RAW").Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}

	return Note{}
}

func GetNotes() []Note {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var notes []Note

	for _, row := range resp.Values {
		notes = append(notes, Note{
			Id:         row[0].(string),
			Title:      row[1].(string),
			Content:    row[2].(string),
			RemindDate: row[3].(string),
			CreatedAt:  row[4].(string),
			UpdatedAt:  row[5].(string),
		})
	}

	return notes
}

func getLastId() int {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return len(resp.Values)
}
