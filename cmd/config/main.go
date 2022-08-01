package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

type Cluster struct {
	Name    string   `yaml:"name"`
	Host    string   `yaml:"host"`
	Port    int      `yaml:"port"`
	User    string   `yaml:"user"`
	Pass    string   `yaml:"pass"`
	Sources []Source `yaml:"sources"`
}

type Source struct {
	GitRepo   string   `yaml:"git-repo"`
	GitBranch string   `yaml:"git-branch"`
	GitPaths  []string `yaml:"git-paths"`
}

func ReadConfigFile(configPath string) (Config, error) {
	var config Config

	// Open YAML file
	file, err := os.Open(configPath)
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
