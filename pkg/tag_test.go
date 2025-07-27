package release

import "testing"

func TestCreateVersionTag(t *testing.T) {
	version := Version{
		Major: 3,
		Minor: 29,
		Patch: 4,
	}
	tagFormat := "version {version}"

	expected := "version 3.29.4"

	result := createVersionTag(version, tagFormat)

	if expected != result {
		t.Errorf("Tag does not match expected output.\nExpected:\n%s\nGot:\n%s", expected, result)
	}
}

func TestParseVersionFromTag(t *testing.T) {
	tag := "v8.33.449"
	tagFormat := "v{version}"

	expected := Version{
		Major: 8,
		Minor: 33,
		Patch: 449,
	}

	result, err := parseVersionFromTag(tag, tagFormat)
	if err != nil {
		t.Fatalf("failed to parse version tag: %v", err)
	}

	if result != expected {
		t.Errorf("Version does not match expected output.\nExpected:\n%v\nGot:\n%v", expected, result)
	}
}

func TestParseVersionString(t *testing.T) {
	versionString := "4.23.5544"

	expected := Version{
		Major: 4,
		Minor: 23,
		Patch: 5544,
	}

	result, err := parseVersionString(versionString)
	if err != nil {
		t.Fatalf("failed to parse version string: %v", err)
	}

	if result != expected {
		t.Errorf("Versiondoes not match expected output.\nExpected:\n%v\nGot:\n%v", expected, result)
	}
}
