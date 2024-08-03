package readfile

import (
	"log"
	"os"
)

func ReadFile(filepath string) (string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	fileString := string(file)
	return fileString, nil
}
