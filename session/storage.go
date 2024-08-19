package session

import (
	"encoding/json"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"log"
	"os"
)

const sessionsFile = "sessions.json"

func saveSessions() {
	file, err := os.Create(sessionsFile)

	if err != nil {
		log.Fatal("Error creating sessions json file", err)
	}

	defer utils.CloseFileError(file)

	var x []Session

	for _, value := range sessions {
		x = append(x, value)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err1 := encoder.Encode(x)

	if err1 != nil {
		log.Fatal("Error encoding sessions json file", err)
	}
}

func ReadSessions() {
	if _, err := os.Stat(sessionsFile); os.IsNotExist(err) {
		return
	}

	file, err := os.Open(sessionsFile)

	if err != nil {
		log.Fatal("Error opening sessions json file", err)
	}

	defer utils.CloseFileError(file)

	var x []Session

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&x)

	if err != nil {
		log.Fatal("Error decoding sessions json file", err)
	}

	for _, value := range x {
		sessions[value.Username] = value
	}

	log.Printf("Read %d sessions.", len(sessions))
}
