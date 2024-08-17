package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name                  string `yaml:"name"`
	BaseDirectory         string `yaml:"baseDirectory"`
	TemplateFile          string `yaml:"templateFile"`
	EnabledArchiveSummary bool   `yaml:"enabledArchiveSummary"`
}

func Load(s string) (*Config, error) {
	config := &Config{}

	err := yaml.UnmarshalStrict([]byte(s), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadFile(filePath string) (*Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return Load(string(content))
}
