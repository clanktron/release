package release

import (
	"fmt"
	"log"
	git "github.com/go-git/go-git/v6"
)

const RELEASE_BRANCH = "main"
const TAG_FORMAT = "{version}"

func Release()  {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}
	tag, commitsSinceRelease := getLatestRelease(getHead(repo, RELEASE_BRANCH), buildTagCommitMap(repo, TAG_FORMAT))
	currentVersion := parseVersionFromTag(tag, TAG_FORMAT)
	fmt.Print(generateChangelog(commitsSinceRelease))
	newVersion := updateVersion(currentVersion, parseSemanticReleaseChangeType(commitsSinceRelease))
	fmt.Println(newVersion.String())
	// update version files/run external program
	// commit changes and tag new version
}
