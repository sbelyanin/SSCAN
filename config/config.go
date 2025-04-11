// config/config.go
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server       ServerConfig `yaml:"server"`
	Logger       LoggerConfig `yaml:"logger"`
	Scans        []ScanConfig `yaml:"scans"`
	ConfigPeriod int          `yaml:"config_period"`
	AuthHashFile string       `yaml:"auth_hash_file"`
}

type ScanConfig struct {
	Path   string `yaml:"path"`
	Period int    `yaml:"period"`
}

type ServerConfig struct {
	Addr    string `yaml:"addr"`
	TLSCert string `yaml:"tls_cert"`
	TLSKey  string `yaml:"tls_key"`
}

type LoggerConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	// Проверка обязательных полей
	if cfg.ConfigPeriod <= 0 {
		cfg.ConfigPeriod = 3600 // Default 1 hour
	}

	return cfg, nil
}
