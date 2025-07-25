package release

import (
	"strconv"
	"strings"
)

func validTagFormat(tag string, tagFormat string) bool {
	// TODO:
	return true
}

func parseVersionFromTag(tag string, tagFormat string) (version Version, err error) {
	// TODO: use tagFormat
	versionComponents := strings.Split(tag, ".")
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

func parseVersionTag(version Version, tagFormat string) string {
	// TODO: use tagFormat
	return version.String()
}
