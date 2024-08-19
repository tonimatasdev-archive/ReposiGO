package utils

import (
	"log"
	"os"
)

func CloseFileError(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Println("Error closing file", file.Name(), ":", err)
	}
}
