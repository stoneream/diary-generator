package data

import "gopkg.in/yaml.v2"

type Metadata struct {
	Title string `yaml:"title"`
	Date  string `yaml:"date"`
}

func (m *Metadata) String() (string, error) {
	data, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
