//go:build integration

package release

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func setupRepo(t *testing.T) (*git.Repository, string) {
	dir, err := os.MkdirTemp("", "test-repo-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}
	return repo, dir
}

func makeTestCommit(t *testing.T, repo *git.Repository, dir string) (*object.Commit){
	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}
	file := filepath.Join(dir, "file.txt")
	err = os.WriteFile(file, []byte("initial content"), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	_, err = wt.Add("file.txt")
	if err != nil {
		t.Fatalf("failed to add file: %v", err)
	}
	commitHash, err := wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Tester",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}
	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		t.Fatalf("failed to get commit object after committing: %v", err)
	}
	return commit
}

func TestWorkingTreeClean(t *testing.T) {
    repo, dir := setupRepo(t)
	makeTestCommit(t, repo, dir)

    t.Run("clean tree", func(t *testing.T) {
        err := workingTreeClean(repo)
        if err != nil {
            t.Errorf("expected clean tree, got error: %v", err)
        }
    })

    t.Run("dirty tree", func(t *testing.T) {
        err := os.WriteFile(filepath.Join(dir, "dirty.txt"), []byte("dirty"), 0644)
        if err != nil {
            t.Fatalf("failed to write dirty file: %v", err)
        }

        err = workingTreeClean(repo)
        if err == nil {
            t.Errorf("expected error for dirty tree, got nil")
        }
    })
}

func getLatestCommit(t *testing.T, repo *git.Repository) *object.Commit {
	ref, err := repo.Head()
	if err != nil {
		t.Fatalf("failed to get HEAD: %v", err)
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		t.Fatalf("failed to get commit object: %v", err)
	}
	return commit
}

func TestGetCurrentHead(t *testing.T) {
	repo, dir := setupRepo(t)
	expectedCommit := makeTestCommit(t, repo, dir)
	commit, err := getCurrentHead(repo, "master")
	if err != nil {
		t.Errorf("unexpected error getting HEAD: %v", err)
	}
	if commit == nil {
		t.Errorf("expected commit, got nil")
	}
	if commit.Hash != expectedCommit.Hash {
		t.Errorf("Commit hash does not match expected output.\nExpected:\n%v\nGot:\n%v", expectedCommit.Hash, commit.Hash)
	}
}

func TestValidateTag(t *testing.T) {
	repo, dir := setupRepo(t)
	makeTestCommit(t, repo, dir)
	tag := "v0.1.0"

	t.Run("valid tag", func(t *testing.T) {
		err := validateTag(repo, tag)
		if err != nil {
			t.Errorf("expected nil for nonexistent tag, got: %v", err)
		}
	})

	t.Run("invalid tag", func(t *testing.T) {
		commit := getLatestCommit(t, repo)
		_, err := repo.CreateTag(tag, commit.Hash, nil)
		if err != nil {
			t.Fatalf("failed to create tag: %v", err)
		}

		err = validateTag(repo, tag)
		if err == nil || err.Error() != "tag already exists" {
			t.Errorf("expected 'tag already exists', got: %v", err)
		}
	})
}

func TestCreateRelease(t *testing.T) {
	repo, _ := setupRepo(t)
	gitConfig := GitConfig{
		Author: "ReleaseBot",
		Email:  "release@example.com",
	}
	tag := "v1.0.0"

	err := CreateRelease(repo, tag, gitConfig)
	if err != nil {
		t.Fatalf("CreateRelease failed: %v", err)
	}
	ref, err := repo.Tag(tag)
	if err != nil {
		t.Fatalf("expected tag to exist, got error: %v", err)
	}
	if ref.Name().Short() != tag {
		t.Errorf("expected tag name %s, got %s", tag, ref.Name().Short())
	}
}

func TestGetLatestRelease(t *testing.T) {
	repo, dir := setupRepo(t)

	wt, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	// Add a few commits
	for i := 0; i < 2; i++ {
		file := filepath.Join(dir, "file.txt")
		err = os.WriteFile(file, []byte(time.Now().String()), 0644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
		_, err = wt.Add("file.txt")
		if err != nil {
			t.Fatalf("failed to add file: %v", err)
		}
		_, err = wt.Commit("update", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "Author",
				Email: "author@example.com",
				When:  time.Now(),
			},
		})
		if err != nil {
			t.Fatalf("commit failed: %v", err)
		}
	}

	// Tag HEAD
	head := getLatestCommit(t, repo)
	_, err = repo.CreateTag("v0.2.0", head.Hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "ReleaseBot",
			Email: "release@example.com",
			When:  time.Now(),
		},
		Message: "Release v0.2.0",
	})
	if err != nil {
		t.Fatalf("failed to tag: %v", err)
	}

	version, commits := getLatestRelease(repo, head, "v{version}")
	if version.String() != "0.2.0" {
		t.Errorf("expected version 0.2.0, got %s", version.String())
	}
	if len(commits) == 0 {
		t.Errorf("expected at least one commit, got 0")
	}
}
