package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"mime/multipart"

	"path/filepath"
)

// CheckFileIsValid checks if the file is valid
func CheckFileIsValid(file *multipart.FileHeader) error {
	extension := filepath.Ext(file.Filename)
	if extension != ".csv" {
		err := errors.New("Unsupported Media type")
		return err
	}
	return nil
}

// ReadFile reads the file
func ReadFile(buf *bytes.Buffer) ([][]string, error) {
	stringReader := buf.String()
	reader := csv.NewReader(bytes.NewBufferString(stringReader))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}
