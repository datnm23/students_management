package writestruct

import (
	"fmt"
	"log"
	"os"
)

func WriteStruct(filepath string, content any) error {
	data := fmt.Sprint(content)
	err := os.WriteFile(filepath, []byte(data), 0644)

	if err != nil {
		log.Fatal(err)
	}
	return nil
}
