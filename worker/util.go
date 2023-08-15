package worker

import (
	"fmt"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/provision"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func TrySSHConnection(node *types.UserNode) (*ssh.Client, error) {
	try_count := 0
	for {
		client, err := provision.CreateClient(fmt.Sprintf("%s:22", node.IPv4))
		if err != nil {
			if try_count > 5 {
				return nil, fmt.Errorf("ssh-client connection failed: %v", err)
			} else {
				logrus.WithFields(logrus.Fields{
					"err":       err,
					"try_count": try_count,
				}).Warn("ssh-client failed")
				try_count++
				time.Sleep(10 * time.Second)
			}
		} else {
			return client, nil
		}
	}
}
