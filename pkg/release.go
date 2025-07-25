package release

import (
	"fmt"
	"log"
	"os/exec"

	git "github.com/go-git/go-git/v6"
)

func Release(configFile string) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	config := DefaultConfig
	if configFile != "" {
		log.Printf("sourcing config file %s...", configFile)
		config, err = parseConfigFile(configFile)
		if err != nil {
			log.Fatalf("error sourcing config file: %s", err.Error())
		}
	}

	if err := workingTreeClean(repo); err != nil {
		log.Fatalf("uncommitted changes detected, exiting...")
	}

	log.Println("getting latest release...")

	startingCommit, err := getCurrentHead(repo, config.ReleaseBranch)
	currentVersion, commitsSinceRelease := getLatestRelease(repo, startingCommit, config.TagFormat)
	log.Printf("current version is %s\n", currentVersion)

	changeType := parseSemanticReleaseChangeType(commitsSinceRelease)
	if changeType == noop {
		log.Fatalf("changes since last release are insufficient - cancelling release...")
	}
	newVersion := updateVersion(currentVersion, changeType)
	log.Printf("%s release - updating version to %s\n", changeType.String(), newVersion.String())

	if err := validateTag(repo, newVersion.String()); err != nil {
		log.Fatalf("error validating new tag against repo: %v", err)
	}

	if config.VersionCommand != "" {
		if repoVersionProcedure(config.VersionCommand) != nil {
			log.Fatalf("version increment command failed - exiting...")
		}
	}

	log.Println("creating release commit and tagging...")
	if CreateRelease(repo, newVersion.String(), config.Git) != nil {
		log.Fatalf("failed to properly create release commit/tag - exiting...")
	}

	fmt.Print(generateChangelog(commitsSinceRelease))
}

// update version files/run external program
func repoVersionProcedure(command string) error {
	return exec.Command("sh", "-c", command).Run()
}
