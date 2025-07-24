package release

import (
	"fmt"
	"strconv"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) String() string {
	return fmt.Sprintf("%s.%s.%s", strconv.Itoa(v.Major), strconv.Itoa(v.Minor), strconv.Itoa(v.Patch))
}

func updateVersion(version Version, changeType semverChangeType) Version {
	switch changeType {
	case major:
		version.Major++
	case minor:
		version.Minor++
	case patch:
		version.Patch++
	}
	return version
}
