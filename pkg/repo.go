package release

import (
	"errors"
	"fmt"
	"time"
	"path"
	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func workingTreeClean(repo *git.Repository) error {
    worktree, err := repo.Worktree()
    if err != nil {
        return err
    }
    status, err := worktree.Status()
    if err != nil {
        return err
    }
	if status.IsClean() {
		return nil
	}
    return fmt.Errorf("working tree is unclean")
}

func getCurrentHead(repo *git.Repository, branch string) (*object.Commit, error) {
	refName := plumbing.NewBranchReferenceName(branch)
	ref, err := repo.Reference(refName, true)
	if err != nil {
		return &object.Commit{}, err
	}
	return repo.CommitObject(ref.Hash())
}

func buildTagMap(repo *git.Repository, tagFormat string) map[plumbing.Hash]string {
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

func getLatestRelease(repo *git.Repository, head *object.Commit, tagFormat string)  (version Version, childCommits []*object.Commit) {
	tagMap := buildTagMap(repo, tagFormat)
	commitIterator := object.NewCommitPreorderIter(head, nil, nil)
	commitIterator.ForEach(func(c *object.Commit) error {
		childCommits = append(childCommits, c)
	    if tagString, ok := tagMap[c.Hash]; ok {
			currentVersion, err := parseVersionFromTag(tagString, tagFormat)
			if err == nil {
				version = currentVersion
				return errors.New("break iterator - found valid release tag")
			}
	    }
	    return nil
	})
	return version, childCommits
}

// returns nil if the tag does not exist
func validateTag(repo *git.Repository, tag string) error {
	_, err := repo.Reference(plumbing.ReferenceName(path.Join("refs", "tags", tag)), false)
	switch err {
	case plumbing.ErrReferenceNotFound:
		return nil
	case nil:
		return fmt.Errorf("tag already exists")
	default:
		return err
	}
}

func CreateRelease(repo *git.Repository, tag string, gitConfig GitConfig) error {
	w, err := repo.Worktree()
	if err != nil || w == nil {
		return fmt.Errorf("error getting worktree: %v", err)
	}
	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return fmt.Errorf("error staging files: %v", err)
	}
	commitHash, err := createReleaseCommit(w, tag, gitConfig)
	if err != nil {
		return fmt.Errorf("error creating release commit: %v", err)
	}
	_, err = createReleaseTag(repo, tag, commitHash, gitConfig)
	if err != nil {
		return fmt.Errorf("error tagging commit: %v", err)
	}
	return nil
}

func createReleaseCommit(w *git.Worktree, tag string, gitConfig GitConfig) (plumbing.ObjectID, error) {
	commitHash, err := w.Commit(fmt.Sprintf("chore(release): %s", tag), &git.CommitOptions{
		AllowEmptyCommits: true,
		Author: &object.Signature{
			Name:  gitConfig.Author,
			Email: gitConfig.Email,
			When:  time.Now(),
		},
	})
	return commitHash, err
}

func createReleaseTag(repo *git.Repository, tag string, commitHash plumbing.ObjectID, config GitConfig) (*plumbing.Reference, error) {
	ref, err := repo.CreateTag(tag, commitHash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  config.Author,
			Email: config.Email,
			When:  time.Now(),
		},
		Message: "Release " + tag,
	})
	return ref, err
}
