package utils

import (
	"encoding/json"
	"github.com/kataras/iris"

	"log"
	"os"
)

// Configuration structure.
type Configuration struct {
	Port             string
	Env              string
	DbDev            string
	DbPD             string
	DbName           string
	UserCollection   string
	GeneratedLinkKey string
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

// ConfigParser get config from context and return confiruration type.
func ConfigParser(ctx iris.Context) Configuration {
	c := ctx.Values().Get("config")
	config, _ := c.(Configuration)

	return config
}
