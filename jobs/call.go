package jobs

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/constants"
	"github.com/abr-ooo/hodhod/log"
	"github.com/abr-ooo/hodhod/models"
	"github.com/abr-ooo/hodhod/services"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeCall = "Call"
)

type CallPayload struct {
	LogList []models.Log
}

func NewCallTask(logList []models.Log) (*asynq.Task, error) {
	payload, err := json.Marshal(CallPayload{LogList: logList})
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"log_list": logList,
			"error":    err.Error(),
		}).Debug("Worker.CALL. create payload json for Call failed!\n", err)
		return nil, err
	}
	return asynq.NewTask(TypeCall, payload), nil
}

func HandleCallTask(ctx context.Context, t *asynq.Task) error {
	var p CallPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Log.Debug("Worker.CALL. unmarshal Call Payload failed: log_list: %s", p.LogList)
		return err
	}
	log.Log.Debug("Worker.CALL. Call: log_list: %s", p.LogList)

	for i := 0; i < len(p.LogList); i++ {
		p.LogList[i].Driver = constants.CALL_KAVEH
		p.LogList[i].Status = constants.INQUEUE

	}

	tx := connections.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Log.WithField("err", r).Debug("Worker.CALL. Recover Failed!")
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Log.WithField("err", err).Debug("Worker.CALL. Transaction Have Error!")
		return err
	}

	if err := tx.Create(&p.LogList).Error; err != nil {
		log.Log.WithField("err", err).Debug("Worker.CALL. Insert Bulk Log in DB Failed!")
		tx.Rollback()
		return err
	}

	// newLogList, err := repositories.CreateBulkLog(p.LogList)
	// if err != nil {
	// 	log.WithField("err", err).Debug("Worker. Insert Bulk Log in DB Failed!")
	// 	return err
	// }
	for key, value := range p.LogList {
		if value.Channel == constants.CALL {
			err := services.Call(value.Message, value.ContactInfo)
			if err != nil {
				log.Log.WithField("err", err).Debug("Worker.CALL. Sending Call Failed!\n", err)
				return err
			}

			if err := tx.Model(&models.Log{}).Where("id = ?", value.ID).Update("status", constants.DONE).Error; err != nil {
				log.Log.WithField("err", err).Debug("Worker.CALL. Update Status Failed!")
				tx.Rollback()
				return err
			}

			log.Log.Debug("Jobs. Call %s for number : %s job finish successfully :) total_call : %s", key, value, strconv.Itoa(len(p.LogList)))

		}
	}

	log.Log.Debug("Jobs. Call Worker job finish successfully :)")

	return tx.Commit().Error
}
