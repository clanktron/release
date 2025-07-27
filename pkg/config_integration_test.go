//go:build integration
package release

import (
	"os"
	"testing"
)

func TestReadConfigFile_Success(t *testing.T) {
	content := []byte("releaseBranch: test-branch")
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	data, err := readConfigFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("readConfigFile failed: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("readConfigFile returned incorrect content.\nExpected: %s\nGot: %s", content, data)
	}
}

func TestReadConfigFile_FileDoesNotExist(t *testing.T) {
	_, err := readConfigFile("nonexistent.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestLoadConfig_Success(t *testing.T) {
	content := `
releaseBranch: test-branch
tagFormat: v{version}
git:
  author: ReleaseBot
  email: releasebot@example.com
versionCommand: "incrementVersion"
`
	tmpFile, err := os.CreateTemp("", "loadconfig-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	expected := Config{
		ReleaseBranch: "test-branch",
		TagFormat:     "v{version}",
		Git: GitConfig{
			Author: "ReleaseBot",
			Email:  "releasebot@example.com",
		},
		VersionCommand: "incrementVersion",
	}

	if cfg != expected {
		t.Errorf("LoadConfig returned unexpected result.\nExpected: %+v\nGot: %+v", expected, cfg)
	}
}
