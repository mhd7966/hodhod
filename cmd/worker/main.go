package main

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/mhd7966/hodhod/configs"
	"github.com/mhd7966/hodhod/connections"
	_ "github.com/mhd7966/hodhod/docs"
	"github.com/mhd7966/hodhod/jobs"
	"github.com/mhd7966/hodhod/log"
	"github.com/sirupsen/logrus"
)

func main() {

	configs.SetConfig()
	log.LogInit()
	connections.ConnectRedis()
	connections.ConnectDatabase()
	defer connections.CloseDb()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: configs.Cfg.Redis.Addr},
		asynq.Config{
			Concurrency: 2,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				sentry.ConfigureScope(func(scope *sentry.Scope) {
					scope.SetContext("Task", map[string]interface{}{
						"Payload": task.Payload(),
						"Type":    task.Type(),
					})
				})
				sentry.CaptureException(err)
			}),
		},
	)

	mux := asynq.NewServeMux()
	log.Log.Info("Server Mux Create Succesfully :)")
	mux.HandleFunc(jobs.TypeSMS, jobs.HandleSMSTask)
	mux.HandleFunc(jobs.TypeCall, jobs.HandleCallTask)

	if err := srv.Run(mux); err != nil {
		log.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("There was a problem running the server")
	}
	log.Log.Info("Server Mux Run Succesfully :)")

}
