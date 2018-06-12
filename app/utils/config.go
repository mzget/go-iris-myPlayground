package utils

import (
	"encoding/json"

	"log"
	"os"
)

// Configuration structure.
type Configuration struct {
	Port  string
	Env   string
	DbDev string
	DbPD  string
}

// GetConfig for load configuration file json.
func GetConfig(confPath string) Configuration {
	file, _ := os.Open(confPath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}

	return configuration
}
