package release

import (
	"fmt"
	"log"
	"os"
	"io"
	"strings"
	"os/exec"

	git "github.com/go-git/go-git/v6"
)

func Release() {

	parseFlags()

	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("error sourcing config file: %v", err)
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}
	if err := workingTreeClean(repo); err != nil {
		log.Fatalf("uncommitted changes detected, exiting...")
	}

	log.Println("getting latest release...")
	startingCommit, err := getCurrentHead(repo, config.ReleaseBranch)
	if err != nil {
		log.Fatalf("failed to get head commit, exiting...")
	}
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

	changelog := generateChangelog(commitsSinceRelease)

	if *dryrun {
		log.Printf("dry-run enabled - version procedure and git operations will be skipped, changelog will not be written to disk")
		fmt.Fprint(os.Stderr, changelog)
	} else {
		if repoVersionProcedure(config.VersionCommand) != nil {
			log.Fatalf("version increment command failed - exiting...")
		}
		log.Println("creating release commit and tagging...")
		if CreateRelease(repo, newVersion.String(), config.Git) != nil {
			log.Fatalf("failed to properly create release commit/tag - exiting...")
		}
		file, err := os.Create("changelog.txt")
		if err != nil {
			log.Fatalf("failed to write changelog.txt: %v", err)
		}
		defer file.Close()
		mw := io.MultiWriter(os.Stderr, file)
		io.Copy(mw, strings.NewReader(changelog))
	}

	fmt.Fprint(os.Stderr, "\nNew Version: ")
	fmt.Printf("%s\n", newVersion.String())
}

// update version files/run external program
func repoVersionProcedure(command string) error {
	if command != "" {
		return exec.Command("sh", "-c", command).Run()
	}
	return nil
}
