package provision

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ProvisionModel struct {
	Config *types.Config
	UserID string
}

func Provision(client *ssh.Client, cfg *types.Config, userID string, status func(string, int)) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("could not open session: %v", err)
	}
	defer session.Close()

	status("creating sftp session", 50)
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("could not create sftp client: %v", err)
	}
	defer sftp.Close()

	status("creating directories", 60)
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

	status("templating files", 65)
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
		return fmt.Errorf("could not write setup.sh: %v", err)
	}

	err = core.SCPWriteFile(sftp, Files, "docker-compose.yml", "/provision/docker-compose.yml", 0644, true)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	status("executing setup script", 80)
	_, stderr, err := core.SSHExecute(client, "/provision/setup.sh")
	if err != nil {
		return fmt.Errorf("SSHExecute error: %v, stderr: '%s'", err, string(stderr))
	}

	return nil
}
