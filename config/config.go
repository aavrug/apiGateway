package config

import (
	"encoding/json"
	"io/ioutil"
)

func init() {
	cfg = new(Config)
}

var cfg *Config

// Config struct
type Config map[string]DiscoverConfig

// DiscoverConfig struct
type DiscoverConfig struct {
	Host string `json:"host,omitempty"`
	Port string `json:"port,omitempty"`
}

// GetConfig method
func GetConfig() *Config {
	return cfg
}

// LoadConfig loads
func LoadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return err
	}
	return nil
}