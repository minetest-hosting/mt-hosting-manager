package server_setup

import (
	"bytes"
	"fmt"
	"html/template"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ComposeModel struct {
	MTUIVersion    string
	Hostname       string
	HTTPRouterName string
}

func Setup(client *ssh.Client, node *types.UserNode, server *types.MinetestServer) error {
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

	basedir := fmt.Sprintf("/data/%s", server.ID)

	err = core.SCPMkDir(sftp, basedir)
	if err != nil {
		return err
	}

	routername := strings.ReplaceAll(server.DNSName, "-", "_")
	routername = strings.ReplaceAll(routername, ".", "_")

	m := &ComposeModel{
		MTUIVersion:    server.UIVersion,
		Hostname:       server.DNSName,
		HTTPRouterName: routername,
	}

	t, err := template.New("").ParseFS(Files, "docker-compose.yml")
	if err != nil {
		return fmt.Errorf("error templating docker-compose: %v", err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, m)
	if err != nil {
		return fmt.Errorf("error executing template 'docker-compose.yml': %v", err)
	}

	err = core.SCPWriteBytes(sftp, buf.Bytes(), fmt.Sprintf("%s/docker-compose.yml", basedir), 0644)
	if err != nil {
		return err
	}

	_, _, err = core.SSHExecute(client, fmt.Sprintf("cd %s && docker-compose pull && docker-compose up -d", basedir))
	if err != nil {
		return err
	}

	return nil
}
