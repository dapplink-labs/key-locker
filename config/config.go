package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database *Database `yaml:"database"`
	Server   *Server   `yaml:"server"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Server struct {
	Port  int  `yaml:"port"`
	Debug bool `yaml:"debug"`
}

func LoadConfigFile(filePath string, cfg *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewDecoder(file).Decode(cfg)
}
