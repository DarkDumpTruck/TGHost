package tghost

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	HTTPPort int `toml:"http_port"`

	IncludeDir string `toml:"include_dir"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTPPort: 40123,
		IncludeDir: "/frontend/assets/include",
	}
}

var globalConfig *Config

func init() {
	globalConfig = DefaultConfig()
}

func GetConfig() *Config {
	return globalConfig
}

func LoadConfigFromFile(path string) (error) {
	cfg := DefaultConfig()
	if _, err := os.Stat(path); err != nil {
		return err
	}
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(buf, cfg); err != nil {
		return err
	}
	globalConfig = cfg
	return nil
}
