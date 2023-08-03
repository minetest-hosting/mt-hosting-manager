package worker

import (
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) MetricsCollector() {
	for {
		nodes, err := w.repos.UserNodeRepo.GetAll()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
				"job": "metrics-collector",
			}).Error("usernode getall failed")
		}

		for _, node := range nodes {
			if node.State != types.UserNodeStateRunning {
				// skip non-running nodes
				continue
			}

			metrics, err := core.FetchMetrics(node.IPv4)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":  err,
					"node": node.ID,
					"ipv4": node.IPv4,
				}).Error("metric collection failed")
			}

			node.LoadPercent = metrics.LoadPercent
			node.DiskSize = metrics.DiskSize
			node.DiskUsed = metrics.DiskUsed
			node.MemorySize = metrics.MemorySize
			node.MemoryUsed = metrics.MemoryUsed
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
				}).Error("node metrics update failed")
			}
		}

		time.Sleep(time.Second * 30)
	}
}
