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

	url := fmt.Sprintf("https://%s.%s/ui", server.DNSName, w.cfg.HostingDomainSuffix)
	dir := "/"
	client := mtui.New(url)
	err = client.Login(server.Admin, server.JWTKey)
	if err != nil {
		return fmt.Errorf("login error: %v", err)
	}

	estimated_size, err := client.GetDirectorySize(dir)
	if err != nil {
		return fmt.Errorf("get dir size error: %v", err)
	}

	r, err := client.DownloadZip(dir)
	if err != nil {
		return fmt.Errorf("download zip error: %v", err)
	}

	//TODO: progress callback / write with scp
	fmt.Print(r, estimated_size)

	//TODO
	return nil
}
