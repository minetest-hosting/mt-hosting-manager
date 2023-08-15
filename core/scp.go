package core

import (
	"bytes"
	"embed"
	"fmt"
	"os"

	"github.com/pkg/sftp"
)

func SCPWriteBytes(sftp *sftp.Client, data []byte, target string, mode os.FileMode) error {
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

func SCPWriteFile(sftp *sftp.Client, fs embed.FS, filename, target string, mode os.FileMode) error {
	data, err := fs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return SCPWriteBytes(sftp, data, target, mode)
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
