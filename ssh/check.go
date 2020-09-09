package ssh

import (
	"fmt"
	"strings"
)

// Check checks if a git-import URL is running under environment where it will most likely to succee
func Check(url string) error {
	//check if url needs ssh
	if !strings.HasPrefix(url, "git@") {
		return nil
	}
	// check if we are running under a unix-like shell environment
	// Otherwise, we don't want to check further
	if !RunInUnixShell() {
		return nil
	}
	// do we have a running ssh-agent?
	ag, err := NewAgentClient()
	if err != nil {
		return err
	}
	defer ag.Close()
	// do we have ssh keys registered with ssh-agent
	if ag.HasKeys() {
		return nil
	}
	// try to add default keys from ~/.ssh as best effort
	err = AddDefaultKeys()
	if err != nil {
		return err
	}
	// check again: do we have ssh keys registered with ssh-agent after adding default keys?
	if ag.HasKeys() {
		return nil
	}
	return fmt.Errorf("no ssh keys registered with ssh-agent. Please use 'ssh-add -k' to add your keys and try again")
}
