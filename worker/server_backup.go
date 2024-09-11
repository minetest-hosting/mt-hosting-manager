package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
)

func (w *Worker) ServerBackup(job *types.Job, status StatusCallback) error {
	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return fmt.Errorf("get server error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}

	backup, err := w.repos.BackupRepo.GetByID(*job.BackupID)
	if err != nil {
		return fmt.Errorf("get backup error: %v", err)
	}
	if backup == nil {
		return fmt.Errorf("backup not found")
	}

	url := fmt.Sprintf("https://%s.%s/ui", server.DNSName, w.cfg.HostingDomainSuffix)
	client := mtui.New(url)
	err = client.Login(server.Admin, server.JWTKey)
	if err != nil {
		return fmt.Errorf("login error: %v", err)
	}

	backup_job, err := client.CreateBackupJob(&mtui.CreateBackupJob{
		Type:     mtui.BackupJobTypeSCP,
		Host:     w.cfg.StorageHostname,
		Username: w.cfg.StorageUsername,
		Password: w.cfg.StoragePassword,
		Port:     22,
		Filename: fmt.Sprintf("%s.tar.gz", backup.ID),
	})
	if err != nil {
		return fmt.Errorf("create backup job error: %v", err)
	}

	job.Message = backup_job.Message
	job.Data = []byte(backup_job.ID)
	return ErrJobStillRunning
}
