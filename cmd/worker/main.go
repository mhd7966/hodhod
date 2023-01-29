package main

import (
	"context"

	"github.com/abr-ooo/hodhod/configs"
	"github.com/abr-ooo/hodhod/connections"
	_ "github.com/abr-ooo/hodhod/docs"
	"github.com/abr-ooo/hodhod/jobs"
	"github.com/abr-ooo/hodhod/log"
	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
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
