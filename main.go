package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zapatacomputing/git-import/ssh"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	url := flag.String("url", "", "the url of the git repository to clone")
	dir := flag.String("dir", "", "the location to put the cloned repository")
	branch := flag.String("branch", "", "the repository branch to clone")
	tag := flag.String("tag", "", "the tag to clone")
	revision := flag.String("revision", "", "a specific revision to checkout")
	flag.Parse()

	err := Clone(*url, *dir, *branch, *tag, *revision)
	if err != nil {
		log.WithFields(log.Fields{
			"url":      *url,
			"dir":      *dir,
			"branch":   *branch,
			"tag":      *tag,
			"revision": *revision,
			"error":    err,
		}).Error("Error cloning git repository")
		os.Exit(1)
	}
}

// Clone clones the git repository at the specified url to the given location
// If cloning a tag or branch, using a 1-commit Clone
// Otherwise the whole repo has to be cloned
func Clone(url string, dir string, branch string, tag string, revision string) error {
	err1 := ssh.Check(url)
	if err1 != nil {
		return fmt.Errorf("git-import: unable to import from [%s] due to error : %w", url, err1)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()

	tag_exists := tag != ""
	branch_exists := branch != ""
	revision_exists := revision != ""
	if !tag_exists && !branch_exists && !revision_exists {
		return fmt.Errorf("Please specify either a branch, a tag, or a revision")
	}

	if (tag_exists && branch_exists) || ((tag_exists || branch_exists) && revision_exists) {
		return fmt.Errorf("Please specify only the branch, only the tag, or only the revision.")
	}

	var referenceName plumbing.ReferenceName
	if branch_exists {
		referenceName = plumbing.NewBranchReferenceName(branch)
	}
	// Tags take precedence over branches so even if a branch was previously specified, we override the reference name with a tag
	if tag_exists {
		referenceName = plumbing.NewTagReferenceName(tag)
	}

	cloneDepth := 1

	// If a revision has a commit hash, we can't clone 1 commit
	if revision_exists {
		cloneDepth = 0
	}

	r, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:           url,
		Depth:         cloneDepth,
		Progress:      os.Stdout,
		SingleBranch:  true,
		ReferenceName: referenceName,
	})
	if err != nil {
		return err
	}
	// go-git cannot clone with an arbitrary revision. We'll have to checkout if a revision was passed
	if revision_exists {
		err := checkoutRevision(r, revision)
		if err != nil {
			// on error, we should remove the repo to match the previous behaviour of a failed clone
			os.RemoveAll(dir)
			return err
		}
	}

	return nil
}

func checkoutRevision(repo *git.Repository, revision string) error {
	// ResolveRevision will give us a hash to checkout (if the revision exists)
	h, err := repo.ResolveRevision(plumbing.Revision(revision))
	if err != nil {
		return err
	}
	// Worktree cannot return an error because we cloned with isBare = false
	w, _ := repo.Worktree()
	// Finally, checkout the revision
	err = w.Checkout(&git.CheckoutOptions{
		Hash: *h,
	})
	if err != nil {
		return err
	}
	return nil
}
