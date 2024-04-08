package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetID = "1XySA2jWwkUU4FLqHGFQGfxHpg6LmDI5TTqLMHNYKcJs"
	readRange     = "Sheet1!A:F"
	credentials   = "key.json"
)

func SheetProcess() []byte {
	// Load the Google Sheets API credentials from your JSON file.
	creds, err := os.ReadFile(credentials)
	if err != nil {
		log.Fatalf("Unable to read credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(creds, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to create JWT config: %v", err)
	}

	client := config.Client(context.Background())
	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Google Sheets service: %v", err)
	}

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(context.Background()).Do()

	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	data, _ := json.Marshal(resp.Values)

	return data
}
