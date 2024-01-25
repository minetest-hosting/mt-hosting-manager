package server_setup

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Restore(client *ssh.Client, cfg *types.Config, node *types.UserNode, server *types.MinetestServer) error {
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

	err = PrepareDataDirectory(sftp, cfg, node, server)
	if err != nil {
		return fmt.Errorf("could not prepare data dir: %v", err)
	}

	basedir := GetBaseDir(server)
	restore_file := fmt.Sprintf("%s/restore.sh", basedir)

	// TODO: template vars for restore script
	err = core.SCPTemplateFile(sftp, Files, "restore.sh", restore_file, 0755, true, map[string]any{})
	if err != nil {
		return err
	}

	_, _, err = core.SSHExecute(client, restore_file)
	if err != nil {
		return err
	}

	return nil
}
