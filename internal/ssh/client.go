package client

import (
	"essh/internal/session"
	"fmt"
	"github.com/helloyi/go-sshclient"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
	"net"
	"os"
	"strings"
)

func addAuthKeyToHost(host, user, passwd string) error {
	pubKey, err := returnPublicKey()

	client, err := sshclient.DialWithPasswd(host+":22", user, passwd)
	if err != nil {
		return err
	}

	addKeyCmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", pubKey)
	addKeyOutput, err := client.Cmd(addKeyCmd).Output()
	if err != nil {
		return fmt.Errorf("Error adding key to authorized_keys: %v\n%s", err, string(addKeyOutput))
	}

	defer client.Close()
	return nil

}
func returnPublicKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	pubKey, err := os.ReadFile(homeDir + "/.ssh/id_rsa.pub")
	pubKey = []byte(strings.ReplaceAll(string(pubKey), "\n", ""))

	return string(pubKey), nil
}

func ConnectWithPassword(session session.Session) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: session.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(session.Password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	connectHost := fmt.Sprintf("%s:%d", session.Host, session.Port)
	sshClient, err := ssh.Dial("tcp", connectHost, config)
	if err != nil {
		return nil, err
	}

	return sshClient, nil
}

func SpawnShell(client *ssh.Client) error {
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to set raw terminal: %s", err)
	}
	defer term.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		return fmt.Errorf("failed to wait for session: %s", err)
	}

	return nil
}

func ConnectWithKey(session session.Session) (*sshclient.Client, error) {
	key, err := returnPublicKey()
	if err != nil {
		return nil, err
	}
	client, err := sshclient.DialWithKey(session.Host+":22", session.UserName, key)
	if err != nil {
		return nil, err
	}

	return client, nil
}
