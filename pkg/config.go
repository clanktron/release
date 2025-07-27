package release

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
    ReleaseBranch  string    `yaml:"releaseBranch"`
    TagFormat      string    `yaml:"tagFormat"`
    Git            GitConfig `yaml:"git"`
    VersionCommand string    `yaml:"versionCommand"`
}

type GitConfig struct {
    Author string `yaml:"author"`
    Email  string `yaml:"email"`
}

var DefaultConfig = Config{
	ReleaseBranch: "main",
	TagFormat:     "{version}",
	Git: GitConfig{
		Author: "Release",
		Email:  "release@example.com",
	},
	VersionCommand: "",
}

func LoadConfig(path string) (Config, error) {
	data, err := readConfigFile(path)
	if err != nil {
		return Config{}, err
	}
	return parseConfig(data)
}

func readConfigFile(path string) (bytes []byte, err error) {
	_, err = os.Stat(path)
	if err != nil {
		return bytes, err
	}
	bytes, err = os.ReadFile(path)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func parseConfig(data []byte) (Config, error) {
	config := DefaultConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}
	return config, nil
}
