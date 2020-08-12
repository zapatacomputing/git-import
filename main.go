package main

import (
	"context"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func main() {
	url := flag.String("url", "", "the url of the git repository to clone")
	dir := flag.String("dir", "", "the location to put the cloned repository")
	branch := flag.String("branch", "master", "the repository branch to clone")
	flag.Parse()

	err := Clone(*url, *dir, *branch)
	if err != nil {
		log.WithFields(log.Fields{
			"url":    *url,
			"dir":    *dir,
			"branch": *branch,
			"error":  err,
		}).Error("Error cloning git repository")
		os.Exit(1)
	}
}

// Clone clones the git repository at the specified url to the given location
// Using a 1-commit clone of the given branch
func Clone(url string, dir string, branch string) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()

	_, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:           url,
		Depth:         1,
		Progress:      os.Stdout,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return err
	}

	return nil
}
