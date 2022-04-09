package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
)

func CheckFileIsValid(file *multipart.FileHeader) error {
	extension := filepath.Ext(file.Filename)
	if extension != ".csv" {
		return fmt.Errorf("File extension is not valid. Only .csv is allowed")
	}
	return nil
}
