package ssh

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// AddDefaultKeys  adds RSA or DSA identities to the authentication agent, ssh-agent(1) located unfder ~/.ssh.
// It adds the files ~/.ssh/id_rsa, ~/.ssh/id_dsa , ~/.ssh/id_ecdsa and ~/.ssh/identity.
// If any file requires a passphrase, ssh-add asks for the passphrase from the user.
// In test automation, ssh keys with passphrase should be abstained.
// Also, a ssh-agent must be running and the SSH_AUTH_SOCK environment variable must contain the name of its socket for ssh-add to work.
func AddDefaultKeys() error {
	bin, err := exec.LookPath("ssh-add")
	if err != nil {
		return fmt.Errorf("could not find 'ssh-add' command")
	}
	cmd := exec.Command(bin)
	cmd.Stderr = new(bytes.Buffer)
	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s failed: %v\n%s", strings.Join(cmd.Args, " "), err, cmd.Stderr)
	}
	return nil
}
