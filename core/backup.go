package core

import (
	"context"
	"fmt"
	"io"
	"mt-hosting-manager/types"

	openssl "github.com/Luzifer/go-openssl/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func StreamBackup(passphrase string, src io.Reader, dst io.Writer) (int64, error) {
	r := openssl.NewReader(src, passphrase, openssl.PBKDF2SHA256)
	return io.Copy(dst, r)
}

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
		return err
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
