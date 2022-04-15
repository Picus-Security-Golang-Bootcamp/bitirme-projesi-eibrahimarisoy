package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
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
func ReadFile(buf *bytes.Buffer) ([][]string, error) {
	fmt.Println("Reading file...")
	fmt.Printf("File size: %s\n", *buf)
	stringReader := buf.String()
	reader := csv.NewReader(bytes.NewBufferString(stringReader))
	print("Creating bulk categories 33333333333")
	fmt.Println("Reading file")

	// reader := csv.NewReader(buf)
	// reader.LazyQuotes = true
	// reader.Comma = ','
	// reader.FieldsPerRecord = -1
	// reader.TrimLeadingSpace = true

	// reader = csv.NewReader(buf)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}
	return records, nil
}
