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
	flag.Parse()

	err := Clone(*url, *dir, *branch, *tag)
	if err != nil {
		log.WithFields(log.Fields{
			"url":    *url,
			"dir":    *dir,
			"branch": *branch,
			"tag":    *tag,
			"error":  err,
		}).Error("Error cloning git repository")
		os.Exit(1)
	}
}

// Clone clones the git repository at the specified url to the given location
// Using a 1-commit clone of the given branch
func Clone(url string, dir string, branch string, tag string) error {
	err1 := ssh.Check(url)
	if err1 != nil {
		return fmt.Errorf("git-import: unable to import from [%s] due to error : %w", url, err1)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()

	if tag == "" && branch == "" {
		return fmt.Errorf("Please specify either a branch or a tag.")
	}

	if tag != "" && branch != "" {
		return fmt.Errorf("Please specify only the branch or only the tag.")
	}

	var referenceName plumbing.ReferenceName
	if branch != "" {
		referenceName = plumbing.NewBranchReferenceName(branch)
	}
	// Tags take precedence over branches so even if a branch was previously specified, we override the reference name with a tag
	if tag != "" {
		referenceName = plumbing.NewTagReferenceName(tag)
	}

	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:           url,
		Depth:         1,
		Progress:      os.Stdout,
		SingleBranch:  true,
		ReferenceName: referenceName,
	})
	if err != nil {
		return err
	}

	return nil
}
