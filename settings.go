package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Location string `yaml:"location"`
	} `yaml:"database"`
	Auth struct {
		Jwt struct {
			Secret string `yaml:"secret"`
		} `yaml:"jwt"`
	} `yaml:"auth"`
}

func ReadSettings() *Settings {
	f, err := os.Open("settings.yml")
	if err != nil {
		log.Fatalf("failed while reading settings.yml: %s\n", err)
	}

	defer f.Close()

	var settings Settings
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&settings)
	if err != nil {
		log.Fatalf("failed while decoding settings.yml: %s\n", err)
	}

	return &settings
}
