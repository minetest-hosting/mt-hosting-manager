package main

import (
	"context"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web"
	"mt-hosting-manager/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
	if cfg.EnableWorker {
		logrus.Info("Starting worker")
		w := worker.NewWorker(repos, cfg)
		go w.Run()
	}
	if cfg.EnableDummyWorker {
		logrus.Info("Starting dummy worker")
		go worker.DummyWorker(repos, cfg)
	}

	// create and setup web api
	api := web.NewApi(repos, cfg)
	api.Setup()

	server := &http.Server{Addr: ":8080", Handler: nil}

	go func() {
		logrus.WithFields(logrus.Fields{"port": 8080}).Info("Listening")
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-captureSignal
	server.Shutdown(context.Background())
}
