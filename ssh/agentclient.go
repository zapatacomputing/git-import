package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/agent"
)

// AgentClient  wraps an instance of agent.Agent
type AgentClient struct {
	socket string
	conn   net.Conn
	agent  agent.Agent
}

// NewAgentClient wrapps a new agent client that uses a unix socket specified by environment variable SSH_AUTH_SOCK
func NewAgentClient() (*AgentClient, error) {
	sshAuthSock := strings.TrimSpace(locateAgentSocket())
	if sshAuthSock == "" {
		return nil, fmt.Errorf("please start a new ssh-agent or register one(via SSH_AUTH_SOCK environment variable) and try again")
	}
	conn, err := net.DialTimeout("unix", sshAuthSock, time.Second)
	if err != nil {
		return nil, fmt.Errorf("error connecting to ssh-agent specified by environment variable SSH_AUTH_SOCK[%s]: %w", sshAuthSock, err)
	}
	return &AgentClient{socket: sshAuthSock, conn: conn, agent: agent.NewClient(conn)}, nil
}

// locateAgentSocket returns ssh-agent auth socket defined by environment variable SSH_AUTH_SOCK
func locateAgentSocket() string {
	return os.Getenv("SSH_AUTH_SOCK")
}

// HasKeys checks if an ssh-agent has any keys registered with it or not
func (ag *AgentClient) HasKeys() bool {
	keys, err := ag.agent.List()
	if err != nil {
		log.Errorf("error while getting key list from ssh-agent : %w", err)
		return false
	}
	return len(keys) > 0
}

// Close - attempts to close this client
func (ag *AgentClient) Close() error {
	err := ag.conn.Close()
	if err != nil {
		log.Errorln(err)
	}
	return err
}
