package core

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SCPWriteBytes(sftp *sftp.Client, data []byte, target string, mode os.FileMode, overwrite bool) error {
	if !overwrite {
		fi, _ := sftp.Stat(target)
		if fi != nil {
			// already exists
			return nil
		}
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

func TemplateFile(fs embed.FS, filename string, model any) ([]byte, error) {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %v", filename, err)
	}

	t, err := template.New("tmpl").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("error templating %s: %v", filename, err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.Execute(buf, model)
	if err != nil {
		return nil, fmt.Errorf("error executing template '%s': %v", filename, err)
	}

	return buf.Bytes(), nil
}

func SCPTemplateFile(sftp *sftp.Client, fs embed.FS, filename, target string, mode os.FileMode, overwrite bool, model any) error {
	data, err := TemplateFile(fs, filename, model)
	if err != nil {
		return err
	}

	err = SCPWriteBytes(sftp, data, target, mode, overwrite)
	if err != nil {
		return fmt.Errorf("template error in file '%s': %v", target, err)
	}

	return nil
}

func SCPWriteFile(sftp *sftp.Client, fs embed.FS, filename, target string, mode os.FileMode, overwrite bool) error {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return SCPWriteBytes(sftp, data, target, mode, overwrite)
}

func SCPMkDir(sftp *sftp.Client, dir string) error {
	_, err := sftp.Stat(dir)
	if err != nil {
		err = sftp.Mkdir(dir)
		if err != nil {
			return fmt.Errorf("could not create directory: %v", err)
		}
	}
	return nil
}

func SSHExecute(client *ssh.Client, cmd string) ([]byte, []byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("could not open session: %v", err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Start(cmd)
	if err != nil {
		return nil, nil, fmt.Errorf("start failed: %v", err)
	}

	err = session.Wait()
	if err != nil {
		ex, ok := err.(*ssh.ExitError)
		if ok {
			fmt.Printf("Exit status: %d\n", ex.ExitStatus())
			if ex.ExitStatus() != 0 {
				return nil, nil, fmt.Errorf("exit-status: %d", ex.ExitStatus())
			}
		} else {
			return nil, nil, fmt.Errorf("unknown script execution error: %v", err)
		}
	}

	return stdout.Bytes(), stderr.Bytes(), nil
}
