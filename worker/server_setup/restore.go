package server_setup

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RestoreModel struct {
	BaseDir       string
	ServerShortID string
	Server        *types.MinetestServer
	Backup        *types.Backup
	Config        *types.Config
}

func Restore(client *ssh.Client, cfg *types.Config, node *types.UserNode, server *types.MinetestServer, backup *types.Backup) error {
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

	m := &RestoreModel{
		BaseDir:       basedir,
		ServerShortID: GetShortName(server.ID),
		Server:        server,
		Config:        cfg,
		Backup:        backup,
	}

	// TODO: template vars for restore script
	err = core.SCPTemplateFile(sftp, Files, "restore.sh", restore_file, 0755, true, m)
	if err != nil {
		return err
	}

	_, _, err = core.SSHExecute(client, restore_file)
	if err != nil {
		return err
	}

	return nil
}
