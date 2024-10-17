package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Record struct {
	ReportType       string
	FirstName        string
	LastName         string
	Email            string
	SentDateUTC      string
	Title            string
	Status           string
	QuizScore        string
	Clicked          string
	ManagerFirstName string
	ManagerLastName  string
	ManagerEmail     string
}

// Function to read CSV files and return filtered records
func parseCSV(fileName, reportType, managerEmail string) ([]Record, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %s", err)
	}

	var result []Record

	// Parse header row
	headers := records[0]
	for _, row := range records[1:] {
		if row[len(row)-1] == managerEmail { // Filter by Manager Email
			record := Record{
				ReportType:       reportType,
				FirstName:        row[indexOf(headers, "First Name")],
				LastName:         row[indexOf(headers, "Last Name")],
				Email:            row[indexOf(headers, "Email")],
				SentDateUTC:      row[indexOf(headers, "Sent Date (UTC)")],
				Title:            row[indexOf(headers, "Title")],
				Status:           getColumnValue(headers, row, "Status"),     // F column
				QuizScore:        getColumnValue(headers, row, "Quiz Score"), // Optional
				Clicked:          getColumnValue(headers, row, "Clicked"),    // Optional
				ManagerFirstName: row[indexOf(headers, "Manager First Name")],
				ManagerLastName:  row[indexOf(headers, "Manager Last Name")],
				ManagerEmail:     row[indexOf(headers, "Manager Email")],
			}
			result = append(result, record)
		}
	}

	return result, nil
}

// Helper function to get column value based on the column name
func getColumnValue(headers []string, row []string, columnName string) string {
	index := indexOf(headers, columnName)
	if index == -1 {
		return "" // Return empty if column not found
	}
	return row[index]
}

// Helper function to get the index of a column in the header row
func indexOf(headers []string, columnName string) int {
	for i, v := range headers {
		if strings.TrimSpace(v) == columnName {
			return i
		}
	}
	return -1
}

// Function to write the output CSV file
func writeCSV(outputFile string, records []Record) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	writer.Write([]string{"ReportType", "First Name", "Last Name", "Email", "Sent Date (UTC)", "Title", "Status", "Quiz Score", "Clicked", "Manager First Name", "Manager Last Name", "Manager Email"})

	// Write data rows
	for _, record := range records {
		writer.Write([]string{
			record.ReportType, record.FirstName, record.LastName, record.Email, record.SentDateUTC, record.Title,
			record.Status, record.QuizScore, record.Clicked, record.ManagerFirstName, record.ManagerLastName, record.ManagerEmail,
		})
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <ManagerEmail> <outputFile>", os.Args[0])
	}

	managerEmail := os.Args[1]
	outputFile := os.Args[2]

	csvFiles := []string{
		"LowScoringUsers.csv",
		"UserIncompleteSessions.csv",
		"UserPhishingFailures.csv",
		"UserIncompleteRemediations.csv",
	}

	var allRecords []Record

	// Process each file
	for _, file := range csvFiles {
		reportType := strings.TrimSuffix(file, ".csv")
		records, err := parseCSV(file, reportType, managerEmail)
		if err != nil {
			log.Fatalf("Error processing file %s: %v", file, err)
		}
		allRecords = append(allRecords, records...)
	}

	// Sort by ReportType and First Name
	sort.Slice(allRecords, func(i, j int) bool {
		if allRecords[i].ReportType == allRecords[j].ReportType {
			return allRecords[i].FirstName < allRecords[j].FirstName
		}
		return allRecords[i].ReportType < allRecords[j].ReportType
	})

	// Write output CSV file
	err := writeCSV(outputFile, allRecords)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	log.Printf("Report generated: %s", outputFile)
}
