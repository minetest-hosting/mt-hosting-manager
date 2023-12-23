package main

import (
	"context"
	"mt-hosting-manager/core"
	"mt-hosting-manager/db"
	"mt-hosting-manager/types"
	"mt-hosting-manager/web"
	"mt-hosting-manager/worker"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis/v8"
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

	// redis/ captcha
	captchaExp := 10 * time.Minute
	var captchaStore captcha.Store = captcha.NewMemoryStore(50, captchaExp)
	if cfg.RedisURL != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     cfg.RedisURL,
			Password: "",
			DB:       0,
		})

		captchaStore = core.NewRedisCaptchaStore(rdb, captchaExp)
	}
	captcha.SetCustomStore(captchaStore)

	// worker (optional)
	var w *worker.Worker
	if cfg.EnableWorker {
		logrus.Info("Starting worker")
		w = worker.NewWorker(repos, cfg)
		w.Start()
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
	logrus.Info("Preparing shutdown")
	if w != nil {
		//shut down worker
		logrus.Info("Shutting down worker")
		w.Stop()
	}
	//stop api
	api.Stop()
	time.Sleep(5 * time.Second)
	db_.Close()
	logrus.Info("Shutdown complete")
	server.Shutdown(context.Background())
}
