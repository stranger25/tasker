package service

import (
	"os"
	"tasker/internal/model"

	"log"

	"gopkg.in/yaml.v2"
)

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
