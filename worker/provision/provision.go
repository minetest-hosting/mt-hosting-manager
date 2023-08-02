package provision

import (
	"bytes"
	"fmt"
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

func writeFile(sftp *sftp.Client, filename, target string, mode os.FileMode) error {
	data, err := Files.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	dstFile, err := sftp.Create(target)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	err = sftp.Chmod(target, mode)
	if err != nil {
		return fmt.Errorf("could not chmod file: %v", err)
	}

	return nil
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

	_, err = sftp.Stat("/provision")
	if err != nil {
		err = sftp.Mkdir("/provision")
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}

	err = writeFile(sftp, "setup.sh", "/provision/setup.sh", 0755)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	err = writeFile(sftp, "docker-compose.yml", "/provision/docker-compose.yml", 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Start("/provision/setup.sh")
	if err != nil {
		return fmt.Errorf("start failed: %v", err)
	}

	err = session.Wait()
	if err != nil {
		ex, ok := err.(*ssh.ExitError)
		if ok {
			fmt.Printf("Exit status: %d\n", ex.ExitStatus())
			if ex.ExitStatus() != 0 {
				return fmt.Errorf("exit-status: %d", ex.ExitStatus())
			}
		} else {
			return fmt.Errorf("unknown script execution error: %v", err)
		}
	}

	return nil
}
