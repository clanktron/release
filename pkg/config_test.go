package release

import "testing"

func TestParseConfig(t *testing.T)  {
	configYaml := `releaseBranch: test-branch
tagFormat: v{version}
git:
  author: ReleaseBot
  email: releasebot@example.com
versionCommand: "incrementVersion"`

	var expected = Config{
		ReleaseBranch: "test-branch",
		TagFormat:     "v{version}",
		Git: GitConfig{
			Author: "ReleaseBot",
			Email:  "releasebot@example.com",
		},
		VersionCommand: "incrementVersion",
	}

	result, err := parseConfig([]byte(configYaml))
	if err != nil {
		t.Fatalf("failed to parse example config: %v", err)
	}

	if result != expected {
		t.Errorf("Changelog does not match expected output.\nExpected:\n%v\nGot:\n%v", expected, result)
	}
}
