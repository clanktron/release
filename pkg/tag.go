package release

import (
	"fmt"
	"strconv"
	"strings"
)

func validTagFormat(tag string, tagFormat string) bool {
	// TODO:
	return true
}

func parseVersionFromTag(tag string, tagFormat string) (Version, error) {
	// Find the placeholder "{version}" in the tagFormat
	placeholder := "{version}"
	placeholderIndex := strings.Index(tagFormat, placeholder)
	if placeholderIndex == -1 {
		return Version{}, fmt.Errorf("invalid tagFormat: missing {version} placeholder")
	}

	// Extract the prefix and suffix from the tagFormat
	prefix := tagFormat[:placeholderIndex]
	suffix := tagFormat[placeholderIndex+len(placeholder):]

	// Strip prefix and suffix from tag
	if !strings.HasPrefix(tag, prefix) || !strings.HasSuffix(tag, suffix) {
		return Version{}, fmt.Errorf("tag does not match format")
	}

	// Extract the version string from the tag
	versionString := strings.TrimPrefix(tag, prefix)
	versionString = strings.TrimSuffix(versionString, suffix)

	// Parse the version string
	return parseVersionString(versionString)
}

func parseVersionString(versionString string) (version Version, err error) {
	versionComponents := strings.Split(versionString, ".")
	version.Major, err = strconv.Atoi(versionComponents[0])
	if err != nil {
		return Version{}, err
	}
	version.Minor, err = strconv.Atoi(versionComponents[1])
	if err != nil {
		return Version{}, err
	}
	version.Patch, err = strconv.Atoi(versionComponents[2])
	if err != nil {
		return Version{}, err
	}
	return version, err
}

func createVersionTag(version Version, tagFormat string) string {
	return strings.ReplaceAll(tagFormat, "{version}", version.String())
}
