package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func Read(filePath string) {

	var result map[string]any
	err := json.Unmarshal([]byte(filePath), &result)
	if err != nil {
		log.Fatal(err)
	}

	filePathArg := result["file_path"].(string)
	fmt.Println(filePathArg)

	data, err := os.ReadFile(filePathArg)

	if err != nil {
		log.Fatal(err)
	}

	write, err := os.Stdout.Write(data)
	if err != nil {
		return
	}

}
