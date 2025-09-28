package runtimeconfig

import (
	"os"
	//"fmt"
	"gopkg.in/yaml.v3"
	"log"
)


var App Config

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Paths struct {
		ConfigDir string `yaml:"config_dir"`
		LogsDir   string `yaml:"logs_dir"`
	} `yaml:"paths"`
}







func LoadAppConfig() {


	data, err := os.ReadFile("./app.yaml")
	if err != nil {
		log.Fatal("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &App); err != nil {
		log.Fatal("failed to unmarshal YAML: %w", err)
	}

}



