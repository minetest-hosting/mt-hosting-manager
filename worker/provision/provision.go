package provision

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"net"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func CreateClient(node *types.UserNode) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:22", node.IPv4)
	key_file := os.Getenv("SSH_KEY")
	f, err := os.ReadFile(key_file)
	if err != nil {
		return nil, fmt.Errorf("ssh-key not found: %s", key_file)
	}

	key, err := ssh.ParsePrivateKey(f)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %v", err)
	}

	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		fp := ssh.FingerprintSHA256(key)
		if node.Fingerprint == "" {
			// no fingerprint yet, add and allow
			node.Fingerprint = fp
		}
		if fp != node.Fingerprint {
			// fingerprint mismatch
			return fmt.Errorf("fingerprint mismatch, on record: '%s'", node.Fingerprint)
		}

		return nil
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: hostKeyCallback,
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

	dirs := []string{
		"/etc/docker",
		"/etc/iptables",
		"/provision",
	}
	for _, dir := range dirs {
		err = core.SCPMkDir(sftp, dir)
		if err != nil {
			return err
		}
	}

	err = core.SCPTemplateFile(sftp, Files, "daemon.json", "/etc/docker/daemon.json", 0644, true, nil)
	if err != nil {
		return err
	}

	err = core.SCPTemplateFile(sftp, Files, "rules.v6", "/etc/iptables/rules.v6", 0644, true, nil)
	if err != nil {
		return err
	}

	err = core.SCPWriteFile(sftp, Files, "setup.sh", "/provision/setup.sh", 0755, true)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	err = core.SCPWriteFile(sftp, Files, "docker-compose.yml", "/provision/docker-compose.yml", 0644, true)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	_, _, err = core.SSHExecute(client, "/provision/setup.sh")
	if err != nil {
		return fmt.Errorf("SSHExecute error: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get cwd: %v", err)
	}

	for _, mmdb_file := range []string{"GeoLite2-ASN.mmdb", "GeoLite2-City.mmdb"} {
		p := path.Join(wd, mmdb_file)
		fi, _ := os.Stat(p)
		if fi != nil {
			data, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("could not read '%s': %v", p, err)
			}

			target := fmt.Sprintf("/data/%s", mmdb_file)
			err = core.SCPWriteBytes(sftp, data, target, 0644, true)
			if err != nil {
				return fmt.Errorf("scp write error '%s': %v", target, err)
			}
		}
	}

	return nil
}
