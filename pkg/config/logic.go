package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadLogic(configPath string) (*Logic, error) {
	config := &Logic{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

type Logic struct {
	Queue `yaml:"queue"`

	Port int
}

type Queue struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Stream   string `yaml:"stream"`
	Group    string `yaml:"group"`
	DB       int    `yaml:"db"`
}
