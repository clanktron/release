package release

import (
	"os"
	"strings"
	"release/pkg/conventionalcommit"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ReleaseBranch           string    `yaml:"releaseBranch"`
	TagFormat               string    `yaml:"tagFormat"`
	TagOnly 				bool 	  `yaml:"tagOnly"`
	Git                     GitConfig `yaml:"git"`
	VersionCommand          string    `yaml:"versionCommand"`
	DryRun                  bool      `yaml:"dryRun"`
	Verbose                 bool      `yaml:"verbose"`
	AllowUncleanWorkingTree bool      `yaml:"allowUncleanWorkingTree"`
	CommitMessage conventionalcommit.Config  `yaml:"commitMessage"`
	MinorTypes              []string  `yaml:"minorTypes"`
	PatchTypes              []string  `yaml:"patchTypes"`
}

type GitConfig struct {
	Author string `yaml:"author"`
	Email  string `yaml:"email"`
}

var defaultConfigFiles = map[string]bool{
	".git-release.yaml": true,
	".git-release.yml":  true,
}

var DefaultConfig = Config{
	ReleaseBranch: "main",
	TagFormat:     "{version}",
	TagOnly: true,
	Git: GitConfig{
		Author: "Release",
		Email:  "release@example.com",
	},
	CommitMessage: conventionalcommit.DefaultConfig,
}

// ReplaceVersionPlaceholder replaces the {version} placeholder in the versionCommand with the actual version.
func ReplaceVersionPlaceholder(versionCommand string, version string) string {
	return strings.ReplaceAll(versionCommand, "{version}", version)
}

func LoadConfig(path string) (config Config, file string, err error) {
	// check if any default config files exist
	if path == "" {
		for file := range defaultConfigFiles {
			_, err = os.Stat(file)
			if err == nil {
				path = file
				break
			}
		}
		if path == "" {
			return DefaultConfig, path, nil
		}
	}
	// load specified config
	_, err = os.Stat(path)
	if err != nil {
		return Config{}, path, err
	}
	data, err := readConfigFile(path)
	if err != nil {
		return Config{}, path, err
	}
	config, err = parseConfig(data)
	return config, path, err
}

func readConfigFile(path string) (bytes []byte, err error) {
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
