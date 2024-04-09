package main

import (
	"context"
	"encoding/json"
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

const (
	spreadsheetID = "1XySA2jWwkUU4FLqHGFQGfxHpg6LmDI5TTqLMHNYKcJs"
	readRange     = "Sheet1!A:F"
	credentials   = "key.json"
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

func CreateNote(note Note) {
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
}

func UpdateNote(note Note) {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for i, row := range resp.Values {
		if row[0] == note.Id {
			resp.Values[i] = []interface{}{note.Id, note.Title, note.Content, note.RemindDate, note.CreatedAt, note.UpdatedAt}
		}
	}

	valueRange := &sheets.ValueRange{
		Values: resp.Values,
	}

	_, err = sheetsService.Spreadsheets.Values.Update(spreadsheetID, readRange, valueRange).ValueInputOption("RAW").Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to write data to sheet: %v", err)
	}
}

func DeleteNote(id string) {
	sheetsService := getClient()

	_, err := sheetsService.Spreadsheets.Values.Clear(spreadsheetID, readRange, &sheets.ClearValuesRequest{}).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to delete data from sheet: %v", err)
	}
}

func getLastId() int {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	return len(resp.Values)
}

func SheetProcess() []byte {
	sheetsService := getClient()

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	data, _ := json.Marshal(resp.Values)

	return data
}
