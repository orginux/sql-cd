package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Clusters []Clusters `yaml:"clusters"`
}

type Clusters struct {
	Name    string   `yaml:"name"`
	Host    string   `yaml:"host"`
	Port    int      `yaml:"port"`
	User    string   `yaml:"user"`
	Pass    string   `yaml:"pass"`
	Sources []Source `yaml:"sources"`
}

type Source struct {
	GitRepo  string   `yaml:"git-repo"`
	GitPaths []string `yaml:"git-paths"`
}

const CONFIG_PATH = "./tests/config.yml"

func ReadConfig() (Config, error) {
	var config Config

	// Open YAML file
	file, err := os.Open(CONFIG_PATH)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	// Decode YAML file to struct
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
