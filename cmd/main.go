package main

import (
	"workflow/cmd/scheduler"
	"workflow/cmd/worker"
	"workflow/config"
	"workflow/src/database"
	"workflow/src/handler"

	"github.com/huylqvn/httpserver"
	"github.com/huylqvn/httpserver/log"
)

// start workflow api server
func main() {
	log := log.Get()

	cfg := config.Load()
	router := httpserver.NewChiRouter(cfg.Port)

	db, err := database.New(cfg)
	if err != nil {
		panic(err)
	}
	log.Info("database connected")

	// close database when server stop
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// init task queue in memory, can alternative use redis, memcached or message queue to store task
	task := make(chan worker.Task)

	if cfg.IsWorker() {
		// worker can become a service worker to process task like send email, push notification, auto scale and use different resource

		log.Info("start worker server")
		workerDispatcher := worker.NewDispatcher(task, cfg.NumberOfWorkers)
		workerDispatcher.Run(db)
		defer workerDispatcher.Stop()
	}

	if cfg.IsScheduler() {
		// scheduler can become a service scheduler for scale and use different resource

		log.Info("start scheduler server")
		job := scheduler.NewScheduler(task, cfg.BatchSize, db)
		job.
			TimeDelay().
			PushScheduler().
			Start()

		defer job.Stop()
	}

	if cfg.IsAPIServer() {
		// api server can become a service api server for scale and use different resource

		log.Info("start workflow api server")
		handler.LoadRouter(router, cfg, db)
	}

	router.
		AddPrefix("/" + cfg.Version).
		AllowCors().
		AllowHealthCheck().
		AllowRecovery().
		ServeHTTP()
}
