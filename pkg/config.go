package release

import (
	"gopkg.in/yaml.v3"
	"log"
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

func parseConfigFile(configFile string) Config {
	config := DefaultConfig
	if configFile == "" {
		log.Printf("no config file - using defaults...")
		return config
	}
	_, err := os.Stat(configFile)
	if err != nil {
		log.Fatalf("failed to parse config file")
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to parse config file")
	}
	if yaml.Unmarshal(configFileBytes, &config) != nil {
		log.Fatalf("failed to parse config file")
	}
	return config
}
