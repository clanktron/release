# git-release

A CLI utility for generating [semantically versioned](https://semver.org) releases based on [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0).

`git-release` is language agnostic as it focuses solely on git operations. It creates clean, traceable release commits and tags, without handling build or publishing steps.

## Features

* Finds the latest release tag on the target branch
* Analyzes commits using Conventional Commit messages
* Generates a changelog
* (Optionally) runs a version update script
* Commits and tags a new release

## Configuration

A config file is optional.

```yaml
releaseBranch: "main"         # Branch to base releases on
tagFormat: "{version}"        # e.g., v{version} or release-{version}
git:
  author: "Release"           # Commit author name
  email: "release@example.com" # Commit author email
versionCommand: ""            # Optional script to update version files
```

> NOTE: `versionCommand` will be run before the release commit. Use it to modify files like `VERSION`, `package.json`, etc.

## What is a "Release"?

A release is a standalone commit that:

* Contains only version file changes (if applicable)
* Is tagged with the version string (e.g., `1.3.0`)
* Represents a single semantic version bump
