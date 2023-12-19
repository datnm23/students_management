package readfile

import (
	"log"
	"os"
)

func ReadFile(filepath string) (string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	fileString := string(file)
	return fileString, err
}
