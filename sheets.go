package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	ErrorsSheetName                 = "All Errors"
	UniqueSheetName                 = "Unique Errors"
	PublishedSheetName              = "Published Products"
	maxColumnsAllowed               = 18278
	minimumUniqueErrorsSheetColumns = 2
	twoWeeksInNanoseconds           = 1.21e+7
	batchRequestLength              = 10000
	metadatTokenURL                 = "http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token"

	styleWarning = iota
	styleError
)

var (
	styles      map[int]map[string]*sheets.Color
	Scopes      map[string]string
	sheetNames  [3]string
	errorHeader []string
)

func main() {
	ctx := context.Background()
	accessToken := "your-google-oauth-access-token-here"
	t := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	})

	sheetsSvc, err := sheets.NewService(ctx, option.WithTokenSource(t))
	if err != nil {
		panic(err)
	}

	svc := sheets.NewSpreadsheetsService(sheetsSvc)
	spreadsheet, err := svc.Create(&sheets.Spreadsheet{}).Context(ctx).Do()
	if err != nil {
		panic(err)
	}

	fmt.Println(spreadsheet)
}

func buildSpreadsheet(sheetNames [3]string, rowCounts [3]int64, columnCounts [3]int64, title string) *sheets.Spreadsheet {
	worksheets := make([]*sheets.Sheet, len(sheetNames))
	spreadsheetProps := buildSpreadsheetProperties(title + " -  Errors")
	for i, sheetName := range sheetNames {
		worksheets[i] = &sheets.Sheet{
			Properties: buildProperties(sheetName, rowCounts[i], columnCounts[i]),
		}
	}

	return &sheets.Spreadsheet{
		Properties: spreadsheetProps,
		Sheets:     worksheets,
	}
}

func buildSpreadsheetProperties(title string) *sheets.SpreadsheetProperties {
	return &sheets.SpreadsheetProperties{
		Title: title,
	}
}

func buildProperties(sheetName string, rowCount, columnCount int64) *sheets.SheetProperties {
	return &sheets.SheetProperties{
		Title: sheetName,
		GridProperties: &sheets.GridProperties{
			ColumnCount: columnCount,
			RowCount:    rowCount,
		},
	}
}
