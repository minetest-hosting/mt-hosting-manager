package server_setup

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SetupModel struct {
	BaseDir       string
	Hostname      string
	Enginename    string
	ServerShortID string
	Server        *types.MinetestServer
	Config        *types.Config
}

func GetShortName(id string) string {
	parts := strings.Split(id, "-")
	return parts[0]
}

func GetEngineName(server *types.MinetestServer) string {
	return fmt.Sprintf("%s_engine", GetShortName(server.ID))
}

const DataDir = "/data"

func GetBaseDir(server *types.MinetestServer) string {
	return fmt.Sprintf("%s/%s", DataDir, server.ID)
}

func Setup(client *ssh.Client, cfg *types.Config, node *types.UserNode, server *types.MinetestServer) error {
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

	err = core.SCPMkDir(sftp, DataDir)
	if err != nil {
		return err
	}

	basedir := GetBaseDir(server)
	err = core.SCPMkDir(sftp, basedir)
	if err != nil {
		return err
	}

	m := &SetupModel{
		BaseDir:       basedir,
		Hostname:      fmt.Sprintf("%s.%s", server.DNSName, cfg.HostingDomainSuffix),
		Enginename:    GetEngineName(server),
		ServerShortID: GetShortName(server.ID),
		Server:        server,
		Config:        cfg,
	}

	if m.Server.UIVersion == "" {
		// fall back to latest
		m.Server.UIVersion = "latest"
	}

	files := []string{
		"docker-compose.yml",
		"nginx.conf",
	}
	for _, filename := range files {
		err = core.SCPTemplateFile(sftp, Files, filename, fmt.Sprintf("%s/%s", basedir, filename), 0644, true, m)
		if err != nil {
			return err
		}
	}

	world_dir := path.Join(basedir, "world")
	err = core.SCPMkDir(sftp, world_dir)
	if err != nil {
		return err
	}

	err = core.SCPTemplateFile(sftp, Files, "minetest.conf", fmt.Sprintf("%s/%s", world_dir, "minetest.conf"), 0644, false, m)
	if err != nil {
		return err
	}

	www_dir := path.Join(world_dir, "www")
	err = core.SCPMkDir(sftp, www_dir)
	if err != nil {
		return err
	}

	err = core.SCPTemplateFile(sftp, Files, "index.html", fmt.Sprintf("%s/%s", www_dir, "index.html"), 0644, false, m)
	if err != nil {
		return err
	}

	setup_file := fmt.Sprintf("%s/setup.sh", basedir)
	err = core.SCPTemplateFile(sftp, Files, "setup.sh", setup_file, 0755, true, m)
	if err != nil {
		return err
	}

	_, _, err = core.SSHExecute(client, setup_file)
	if err != nil {
		return err
	}

	return nil
}
