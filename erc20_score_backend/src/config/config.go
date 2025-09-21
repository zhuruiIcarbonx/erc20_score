package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

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
	Blocakchain struct {
		Chian1 struct {
			Erc20Address string `yaml:"erc20_address"`
			KeyUrl       string `yaml:"key_url"`
			ChainID      string `yaml:"chain_id"`
		} `yaml:"chian1"`
		Chian2 struct {
			Erc20Address string `yaml:"erc20_address"`
			KeyUrl       string `yaml:"key_url"`
			ChainID      string `yaml:"chain_id"`
		} `yaml:"chian2"`
	} `yaml:"blocakchain"`
}

func GetConfig() Config {

	proj, _ := filepath.Abs(filepath.Dir("."))
	fmt.Println("项目路径： %s", proj)
	// 目标文件路径
	configFile, err := ioutil.ReadFile(proj + "/config/config.yml")
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v", err)
	}

	fmt.Printf("config: %v\n", config)
	return config
}
