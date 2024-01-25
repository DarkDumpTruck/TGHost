package tghost

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	HTTPPort int `toml:"http_port"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTPPort: 40123,
	}
}

func LoadConfigFromFile(path string) (*Config, error) {
	cfg := DefaultConfig()
	if _, err := os.Stat(path); err != nil {
		fmt.Println("Config file not found, using default config")
		return cfg, nil
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
