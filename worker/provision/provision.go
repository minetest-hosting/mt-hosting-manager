package provision

import (
	"fmt"
	"mt-hosting-manager/core"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func CreateClient(addr string) (*ssh.Client, error) {
	key_file := os.Getenv("SSH_KEY")
	f, err := os.ReadFile(key_file)
	if err != nil {
		return nil, fmt.Errorf("ssh-key not found: %s", key_file)
	}

	key, err := ssh.ParsePrivateKey(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		// TODO: remember fingerprint from first connection
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("dial error: %v", err)
	}
	return client, nil
}

func Provision(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("could not open session: %v", err)
	}
	defer session.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("could not create sftp client: %v", err)
	}
	defer sftp.Close()

	err = core.SCPMkDir(sftp, "/provision")
	if err != nil {
		return err
	}

	err = core.SCPWriteFile(sftp, Files, "setup.sh", "/provision/setup.sh", 0755)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	err = core.SCPWriteFile(sftp, Files, "docker-compose.yml", "/provision/docker-compose.yml", 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	_, _, err = core.SSHExecute(client, "/provision/setup.sh")
	if err != nil {
		return fmt.Errorf("SSHExecute error: %v", err)
	}

	return nil
}
