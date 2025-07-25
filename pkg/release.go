package release

import (
	"fmt"
	"log"
	"os/exec"

	git "github.com/go-git/go-git/v6"
)

const DEFAULT_CONFIG_FILE = ".release"

func Release() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	config := parseConfigFile("")
	if !workingTreeClean(repo) {
		log.Fatalf("uncommitted changes detected, exiting...")
	}

	log.Println("getting latest release...")
	tag, commitsSinceRelease := getLatestRelease(getHead(repo, config.ReleaseBranch), buildTagCommitMap(repo, config.TagFormat))
	currentVersion := parseVersionFromTag(tag, config.TagFormat)
	log.Printf("current version is %s\n", currentVersion)

	changeType := parseSemanticReleaseChangeType(commitsSinceRelease)
	newVersion := updateVersion(currentVersion, changeType)

	log.Printf("%s release - updating version to %s\n", changeType.String(), newVersion.String())
	if config.VersionCommand != "" {
		if repoVersionProcedure(config.VersionCommand) != nil {
			log.Fatalf("version increment command failed - exiting...")
		}
	}

	log.Println("creating release commit and tagging...")
	_, err = createReleaseCommit(repo, newVersion, config.TagFormat, config.Git)
	if err != nil {
		log.Fatalf("failed to properly create release commit/tag: %s - exiting...", err.Error())
	}

	fmt.Print(generateChangelog(commitsSinceRelease))
}

// update version files/run external program
func repoVersionProcedure(command string) error {
	return exec.Command("sh", "-c", command).Run()
}
