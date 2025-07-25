package release

import (
	"errors"
	"fmt"
	"log"
	"time"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func workingTreeClean(repo *git.Repository) bool {
    worktree, err := repo.Worktree()
    if err != nil {
        log.Fatalf("failed to get worktree: %v", err)
    }
    status, err := worktree.Status()
    if err != nil {
        log.Fatalf("failed to get status: %v", err)
    }
    return status.IsClean()
}

func getHead(repo *git.Repository, branch string) *object.Commit {
	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := repo.Reference(refName, true)
	if err != nil {
		log.Fatalf("Failed to get reference for branch %s: %v", branch, err)
	}
	headCommit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		log.Fatalf("Failed to get commit object: %v", err)
	}
	return headCommit
}

func buildTagCommitMap(repo *git.Repository, tagFormat string) map[plumbing.Hash]string {
	tagMap := make(map[plumbing.Hash]string)
	tags, _ := repo.Tags()
	tags.ForEach(func(ref *plumbing.Reference) error {
		if validTagFormat(ref.Name().Short(), tagFormat) {
	    	hash := ref.Hash()
	    	tagObj, err := repo.TagObject(hash) // lightweight tag
	    	if err == nil {
	    	    // Annotated tag
	    	    hash = tagObj.Target
	    	}
	    	tagMap[hash] = ref.Name().Short()
		}
	    return nil
	})
	return tagMap
}

func getLatestRelease(head *object.Commit, tagMap map[plumbing.Hash]string)  (tag string, childCommits []*object.Commit) {
	commitIterator := object.NewCommitPreorderIter(head, nil, nil)
	commitIterator.ForEach(func(c *object.Commit) error {
		childCommits = append(childCommits, c)
	    if tagString, ok := tagMap[c.Hash]; ok {
			tag = tagString
			return errors.New("break iterator - found tag")
	    }
	    return nil
	})
	return tag, childCommits
}

func createReleaseCommit(repo *git.Repository, version Version, tagFormat string, gitConfig GitConfig) (*object.Commit, error) {
	w, err := repo.Worktree()
	if err != nil {
		return &object.Commit{}, err
	}
	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return &object.Commit{}, err
	}
	commitHash, err := w.Commit(fmt.Sprintf("chore(release): %s", version.String()), &git.CommitOptions{
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  gitConfig.Author,
			Email: gitConfig.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return &object.Commit{}, err
	}
	// TODO: use version and tagFormat to construct tag
	tagName := version.String()
	_, err = repo.CreateTag(tagName, commitHash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  gitConfig.Author,
			Email: gitConfig.Email,
			When:  time.Now(),
		},
		Message: "Release " + tagName,
	})
	if err != nil {
		return &object.Commit{}, err
	}
	commit, err := repo.CommitObject(commitHash)
	if err != nil {
		return &object.Commit{}, err
	}
	return commit, nil
}
