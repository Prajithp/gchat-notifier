package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Channels []Channel `yaml:"channels"`
}

type Channel struct {
	Name   string   `yaml:"name"`
	Url    string   `yaml:"url"`
	Labels []string `yaml:"labels"`
}

func ReadConfig(config_file string) (*Config, error) {
	file, err := ioutil.ReadFile(config_file)
	if err != nil {
		return nil, err
	}
	c := &Config{}

	err = yaml.Unmarshal(file, c)

	if err != nil {
		return nil, err
	}
	return c, nil
}
