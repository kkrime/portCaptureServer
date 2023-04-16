package config

import (
	"github.com/BurntSushi/toml"
)

type DBConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Dbname   string
}
type Config struct {
	DBConfig    DBConfig
	PortCapture PortCapture
}

type PortCapture struct {
	WorkerThreads int
}

func ReadConfig(configFilePath string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(configFilePath, &config)
	return &config, err
}
