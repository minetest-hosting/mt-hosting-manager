package core

import (
	"context"
	"fmt"
	"mt-hosting-manager/types"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

func (c *Core) GetS3Client() (*minio.Client, error) {
	return minio.New(c.cfg.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.cfg.S3KeyID, c.cfg.S3AccessKey, ""),
		Secure: true,
	})
}

func getBackupFilename(b *types.Backup) string {
	return fmt.Sprintf("backup/%s.tar.gz", b.ID)
}

func (c *Core) RemoveBackup(b *types.Backup) error {
	client, err := c.GetS3Client()
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = client.RemoveObject(ctx, c.cfg.S3Bucket, getBackupFilename(b), minio.RemoveObjectOptions{})
	if err != nil {
		// ignore errors while removing
		logrus.WithFields(logrus.Fields{
			"backup_id": b.ID,
			"error":     err,
		}).Error("removing backup from s3 storage")
	}

	return c.repos.BackupRepo.Delete(b.ID)
}

func (c *Core) GetBackupSize(b *types.Backup) (int64, error) {
	client, err := c.GetS3Client()
	if err != nil {
		return 0, err
	}

	ctx := context.Background()
	info, err := client.StatObject(ctx, c.cfg.S3Bucket, getBackupFilename(b), minio.GetObjectOptions{})
	if err != nil {
		return 0, err
	}
	return info.Size, nil
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

	fmt.Printf("Backup stub %v, %v\n", node, server)
	return nil
}
