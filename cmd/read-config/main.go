package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Clusters []Clusters `yaml:"clusters"`
}

type Clusters struct {
	Name   string   `yaml:"name"`
	Host   string   `yaml:"host"`
	Port   int      `yaml:"port"`
	User   string   `yaml:"user"`
	Pass   string   `yaml:"pass"`
	Source []Source `yaml:"source"`
}

type Source struct {
	GitRepo string `yaml:"git-repo"`
	GitPath string `yaml:"git-path"`
}

const CONFIG_PATH = "./tests/config.yml"

func ReadConfig() Config {
	var config Config

	// Open YAML file
	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	// Decode YAML file to struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Println(err.Error())
	}

	return config
}
