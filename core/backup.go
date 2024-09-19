package core

import (
	"fmt"
	"io"
	"mt-hosting-manager/types"

	"github.com/pkg/sftp"
)

func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("%s.tar.gz", b.ID)
}

func (c *Core) RemoveBackup(b *types.Backup) error {
	client, err := CreateStorageClient(c.cfg)
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}
	defer client.Close()

	sc, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("sftp client error: %v", err)
	}
	defer sc.Close()

	err = sc.Remove(getBackupFilename(b))
	if err != nil {
		return fmt.Errorf("sftp remove error: %v", err)
	}

	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) GetBackupSize(b *types.Backup) (int64, error) {
	client, err := CreateStorageClient(c.cfg)
	if err != nil {
		return 0, fmt.Errorf("create client error: %v", err)
	}
	defer client.Close()

	sc, err := sftp.NewClient(client)
	if err != nil {
		return 0, fmt.Errorf("sftp client error: %v", err)
	}
	defer sc.Close()

	fi, err := sc.Stat(getBackupFilename(b))
	if err != nil {
		return 0, fmt.Errorf("sftp stat error: %v", err)
	}

	return fi.Size(), nil
}

func (c *Core) StreamBackup(b *types.Backup, w io.Writer) error {
	client, err := CreateStorageClient(c.cfg)
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}
	defer client.Close()

	sc, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("sftp client error: %v", err)
	}
	defer sc.Close()

	f, err := sc.Open(getBackupFilename(b))
	if err != nil {
		return fmt.Errorf("scp file open error: %v", err)
	}
	defer f.Close()

	var reader io.Reader
	reader = f

	if b.Passphrase != "" {
		// enable decryption
		reader, err = EncryptedReader(b.Passphrase, f)
		if err != nil {
			return fmt.Errorf("decryption failed: %v", err)
		}
	}

	_, err = io.Copy(w, reader)
	return err
}

func (c *Core) RemoveBackupSpace(bs *types.BackupSpace) error {
	list, err := c.repos.BackupRepo.GetByBackupSpaceID(bs.ID)
	if err != nil {
		return err
	}

	for _, b := range list {
		err = c.RemoveBackup(b)
		if err != nil {
			return err
		}
	}

	return c.repos.BackupSpaceRepo.Delete(bs.ID)
}
