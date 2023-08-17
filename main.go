package main

import (
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web"
	"mt-hosting-manager/worker"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if os.Getenv("LOGLEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	db_, err := db.Init(wd)
	if err != nil {
		panic(err)
	}

	err = db.Migrate(db_)
	if err != nil {
		panic(err)
	}

	cfg := types.NewConfig()
	repos := db.NewRepositories(db_)

	// worker (optional)
	if os.Getenv("ENABLE_WORKER") == "true" {
		logrus.Info("Starting worker")
		w := worker.NewWorker(repos, cfg)
		go w.Run()
	}

	// web (always on)
	logrus.WithFields(logrus.Fields{
		"port": 8080,
	}).Info("Starting webserver")

	err = web.Serve(repos, cfg)
	if err != nil {
		panic(err)
	}
}
