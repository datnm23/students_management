package writestruct

import (
	"fmt"
	"log"
	"os"
)

func WriteStruct(filepath string, content any) error {
	data := fmt.Sprint(content)
	return os.WriteFile(filepath, []byte(data), 0644)
}
