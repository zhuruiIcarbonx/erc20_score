package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port     int `yaml:"port"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Userpassword struct {
		Key string `yaml:"key"`
	} `yaml:"userpassword"`
}

func GetConfig() Config {
	fmt.Println("---------------------------")
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("Error parsing YAML: %v", err)
	}

	fmt.Printf("config: %v\n", config)
	return config
}
