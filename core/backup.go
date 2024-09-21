package core

import (
	"fmt"
	"io"
	"mt-hosting-manager/types"
)

func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("%s.tar.gz", b.ID)
}

func (c *Core) RemoveBackup(b *types.Backup) error {
	client, err := CreateStorageClient(c.cfg)
	if err != nil {
		return fmt.Errorf("create client error: %v", err)
	}

	err = client.Remove(getBackupFilename(b))
	if err != nil {
		return fmt.Errorf("webdav remove error: %v", err)
	}

	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) GetBackupSize(b *types.Backup) (int64, error) {
	client, err := CreateStorageClient(c.cfg)
	if err != nil {
		return 0, fmt.Errorf("create client error: %v", err)
	}

	fi, err := client.Stat(getBackupFilename(b))
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

	r, err := client.ReadStream(getBackupFilename(b))
	if err != nil {
		return fmt.Errorf("readstream error: %v", err)
	}
	defer r.Close()

	var reader io.Reader
	reader = r

	if b.Passphrase != "" {
		// enable decryption
		reader, err = EncryptedReader(b.Passphrase, reader)
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
