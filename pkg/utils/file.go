package utils

import (
	"encoding/csv"
	"errors"
	"mime/multipart"

	"path/filepath"
)

// CheckFileIsValid checks if the file is valid
func CheckFileIsValid(file *multipart.FileHeader) error {
	extension := filepath.Ext(file.Filename)
	if extension != ".csv" {
		err := errors.New("File extension is not .csv.")
		return err
	}
	return nil
}

// ReadFile reads the file
func ReadFile(filename *multipart.FileHeader) ([][]string, error) {
	file, err := filename.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.TrimLeadingSpace = true

	reader = csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
