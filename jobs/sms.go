package jobs

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/mhd7966/hodhod/configs"
	"github.com/mhd7966/hodhod/connections"
	"github.com/mhd7966/hodhod/constants"
	"github.com/mhd7966/hodhod/log"
	"github.com/mhd7966/hodhod/models"
	"github.com/mhd7966/hodhod/services"
	"github.com/sirupsen/logrus"
)

const (
	TypeSMS = "SMS"
)

type SMSPayload struct {
	Message string
	Numbers []string
	LogList []models.Log
}

func NewSMSTask(message string, numbers []string, logList []models.Log) (*asynq.Task, error) {
	payload, err := json.Marshal(SMSPayload{Message: message, Numbers: numbers, LogList: logList})
	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"message":  message,
			"numbers":  numbers,
			"log_list": logList,
			"error":    err.Error(),
		}).Debug("Worker.SMS. create payload json for SMS failed!\n", err)
		return nil, err
	}
	return asynq.NewTask(TypeSMS, payload), nil
}

func HandleSMSTask(ctx context.Context, t *asynq.Task) error {
	var p SMSPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Log.Debug("Worker.SMS. unmarshal SMS Payload failed: message=%s and numbers=%s!", p.Message, p.Numbers)
		return err
	}
	log.Log.Debug("Worker.SMS. SMS: message=%s and numbers=%s!", p.Message, p.Numbers)

	tx := connections.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Log.WithField("err", r).Debug("Worker.SMS. Recover Failed!")
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Log.WithField("err", err).Debug("Worker.SMS. Transaction Have Error!")
		return err
	}

	c := configs.Cfg.SMS
	driver, err := strconv.Atoi(c.Driver)
	if err != nil {
		return err
	}

	switch driver {
	case constants.SMS_KAVEH:
		for i := 0; i < len(p.LogList); i++ {
			p.LogList[i].Driver = constants.SMS_KAVEH
			p.LogList[i].Status = constants.INQUEUE
		}
		if err := tx.Create(p.LogList).Error; err != nil {
			log.Log.WithField("err", err).Debug("Worker.SMS. Insert Bulk Log in DB Failed!")
			tx.Rollback()
			return err
		}
		services.KavenegarSendSMS(p.Message, p.Numbers)

	case constants.SMS_SIGNAL:
		for i := 0; i < len(p.LogList); i++ {
			p.LogList[i].Driver = constants.SMS_KAVEH
			p.LogList[i].Status = constants.INQUEUE
		}
		if err := tx.Create(p.LogList).Error; err != nil {
			log.Log.WithField("err", err).Debug("Worker.SMS. Insert Bulk Log in DB Failed!")
			tx.Rollback()
			return err
		}
		services.SignalSendSMS(p.Message, p.Numbers)

	default:
		for i := 0; i < len(p.LogList); i++ {
			p.LogList[i].Driver = constants.SMS_KAVEH
			p.LogList[i].Status = constants.INQUEUE
		}
		if err := tx.Create(p.LogList).Error; err != nil {
			log.Log.WithField("err", err).Debug("Worker.SMS. Insert Bulk Log in DB Failed!")
			tx.Rollback()
			return err
		}
		services.SignalSendSMS(p.Message, p.Numbers)
	}

	// if err := tx.Create(&p.LogList).Error; err != nil {
	// 	log.WithField("err", err).Debug("Worker. Insert Bulk Log in DB Failed!")
	// 	tx.Rollback()
	// 	return err
	// }

	// newLogList, err := services.SendSMS(tx, p.Message, p.Numbers, p.LogList)
	// if err != nil {
	// 	log.Debug("Worker. Sending SMS Failed!\n", err)
	// 	return err
	// }

	for _, value := range p.LogList {
		if err := tx.Model(&models.Log{}).Where("id = ?", value.ID).Update("status", constants.DONE).Error; err != nil {
			log.Log.WithField("err", err).Debug("Worker.SMS. Update Status Failed!")
			tx.Rollback()
			return err
		}
	}

	log.Log.Debug("Jobs.SMS Worker job finish successfully :)")

	return tx.Commit().Error
}
