package release

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ReleaseBranch  string
	TagFormat      string
	Git            GitConfig
	VersionCommand string
}

type GitConfig struct {
	Author string
	Email  string
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

func parseConfigFile(configFile string) (Config, error) {
	config := DefaultConfig
	_, err := os.Stat(configFile)
	if err != nil {
		return config, err
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(configFileBytes, &config); err != nil {
		return config, err
	}
	return config, nil
}
