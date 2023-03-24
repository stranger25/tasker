package service

import (
	"os"
	"tasker/internal/model"

	"log"

	"gopkg.in/yaml.v2"
)

// InitConfig opens the config file and reads the configuration from it.
// It returns a pointer to the Config struct and an error if one is encountered.
// It also prints out log messages to indicate the success of the operations.
func InitConfig() (cfg *model.Config, err error) {
	file, err := os.Open("../config/config.yaml")

	if err != nil {
		return cfg, err
	}
	log.Println("Load config file .. ok")

	decoder := yaml.NewDecoder(file)

	var config model.Config
	err = decoder.Decode(&config)

	if err != nil {
		return &config, err
	}

	log.Println("Read config file .. ok")
	return &config, nil
}
