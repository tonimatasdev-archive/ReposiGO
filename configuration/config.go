package configuration

import (
	"encoding/json"
	"github.com/TonimatasDEV/ReposiGO/repo"
	"github.com/TonimatasDEV/ReposiGO/utils"
	"log"
	"os"
)

type Config struct {
	Port         int          `json:"port"`
	Primary      string       `json:"primaryRepository"`
	CertFile     string       `json:"certFile"`
	KeyFile      string       `json:"keyFile"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

const configFile = "config.json"

func LoadConfig() (*Config, error) {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Println("Config file not found, creating default configuration...")
		defaultConfig := Config{
			Port:     8080,
			Primary:  "releases",
			CertFile: "none",
			KeyFile:  "none",
			Repositories: []Repository{
				{"Releases", "releases", repo.Public},
				{"Secret", "secret", repo.Secret},
				{"Private", "private", repo.Private}},
		}

		file, err := os.Create(configFile)

		if err != nil {
			return nil, err
		}

		defer utils.FileError(file)

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(defaultConfig); err != nil {
			return nil, err
		}

		return &defaultConfig, nil
	}

	file, err := os.Open(configFile)

	if err != nil {
		return nil, err
	}

	defer utils.FileError(file)

	var config Config
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
