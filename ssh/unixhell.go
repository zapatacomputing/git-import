package ssh

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// RunInUnixShell determines if git-import is being run in a unix shell environment or not
// Note: valid shell environments are:
//        MacOS, Linux, Git Bash, Cygwin, MingX, WSL
func RunInUnixShell() bool {
	bin, err := exec.LookPath("sh")
	if err != nil {
		log.Errorf("could not find 'sh' command. Probably running in a non-Unix-like shell environment")
		return false
	}
	cmd := exec.Command(bin, "-c", "uname", "-s")
	cmd.Env = []string{} // Do not let the user's environment influence ssh-agent behavior.
	cmd.Stderr = new(bytes.Buffer)
	out, err := cmd.Output()
	if err != nil {
		log.Errorf("%s failed: %v\n%s", strings.Join(cmd.Args, " "), err, cmd.Stderr)
		return false
	}
	log.Infoln("running in Unix-like shell environment : ", string(out))
	return true
}
