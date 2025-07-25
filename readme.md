# git-release

A cli utility for using [conventional commit messages](https://www.conventionalcommits.org/en/v1.0.0) to auto-create [semantically versioned](https://semver.org) releases.

>Releases are defined here as standalone commits that exclusively contain changes to version files. These commits are to be tagged with their corresponding version.

`git-release` will:
- find the latest release (if it exists)
- generate a changelog for commits following said release
- optionally run an external command to update version files
- commit and tag a new version
