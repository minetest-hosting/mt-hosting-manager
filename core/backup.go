package core

import (
	"fmt"
	"mt-hosting-manager/types"
)

/*
func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("%s.tar.gz", b.ID)
}
*/

func (c *Core) RemoveBackup(b *types.Backup) error {
	// TODO
	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) GetBackupSize(b *types.Backup) (int64, error) {
	// TODO
	return 0, nil
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

func (c *Core) StartBackup(b *types.Backup) error {
	server, err := c.repos.MinetestServerRepo.GetByID(b.MinetestServerID)
	if err != nil {
		return fmt.Errorf("server fetch error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", b.MinetestServerID)
	}

	node, err := c.repos.UserNodeRepo.GetByID(server.UserNodeID)
	if err != nil {
		return fmt.Errorf("usernode fetch error: %v", err)
	}
	if node == nil {
		return fmt.Errorf("usernode not found: %s", server.UserNodeID)
	}

	// TODO
	fmt.Printf("Backup stub %v, %v\n", node, server)
	return nil
}
