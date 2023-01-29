package controllers

import (
	"fmt"
	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/constants"
	"github.com/abr-ooo/hodhod/inputs"
	"github.com/abr-ooo/hodhod/jobs"
	"github.com/abr-ooo/hodhod/models"
	"github.com/abr-ooo/hodhod/outputs"
	repositories "github.com/abr-ooo/hodhod/repositoies"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/abr-ooo/hodhod/log"

)

var logBody models.Log

// Hoo godoc
// @Summary hoo
// @Description hoo
// @ID hoo
// @Tags Hoo
// @Param hoo body inputs.Hoo true "Hoo struct"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /hoo [post]
func Hoo(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	hoo := new(inputs.Hoo)
	err := c.BodyParser(hoo)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("Hoo. Parse body to Hoo failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	logBody = models.Log{
		UserID:    hoo.UserID,
		AlertID:   hoo.AlertID,
		ProjectID: hoo.ProjectID,
		Service:   hoo.Service,
	}

	contacts, err := repositories.GetHooContact(hoo.UserID, hoo.AlertID, hoo.ProjectID, hoo.Service)
	if err != nil {
		response.Message = "Get Contact Info Failed"
		log.Log.WithFields(logrus.Fields{
			"hoo":      hoo,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("Hoo. There is no contact for this request!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alert, err := repositories.GetAlert(hoo.AlertID)
	if err != nil {
		response.Message = "Get Alert Info Failed"
		log.Log.WithFields(logrus.Fields{
			"hoo":      hoo,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("Hoo. There is no alert for this ID!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validContacts := CheckSeverity(*contacts, alert.Severity)
	list, message, logList, err := GetSendingList(validContacts, alert.Message, hoo.Values)
	if err != nil {
		response.Message = "Get Sending List Failed"
		log.Log.WithFields(logrus.Fields{
			"hoo":            hoo,
			"valid_contacts": validContacts,
			"response":       response.Message,
			"error":          err.Error(),
		}).Error("Hoo.Getting info of contacts failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var smsLogList []models.Log
	for _, value := range logList {
		if value.Channel == constants.SMS {
			smsLogList = append(smsLogList, value)
		}
	}

	var callLogList []models.Log
	for _, value := range logList {
		if value.Channel == constants.CALL {
			callLogList = append(callLogList, value)
		}
	}

	if len(smsLogList) != 0 {
		smsTask, err := jobs.NewSMSTask(message, list.SMSList, smsLogList)
		if err != nil {
			response.Message = "Hoo - could not create sms task"
			log.Log.WithFields(logrus.Fields{
				"message": message,
				"numbers": list.SMSList,
				"error":   err.Error(),
			}).Error("Hoo. Couldn't create new SMS task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		_, err = connections.RedisClient.Enqueue(smsTask)
		if err != nil {
			response.Message = "SMS - could not enqueue task"
			log.Log.WithFields(logrus.Fields{
				"message":  message,
				"sms_list": list.SMSList,
				"task":     smsTask,
				"error":    err.Error(),
			}).Error("Hoo. Couldn't enqueue sms task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if len(callLogList) != 0 {
		callTask, err := jobs.NewCallTask(callLogList)
		if err != nil {
			response.Message = "Hoo - could not create call task"
			log.Log.WithFields(logrus.Fields{
				"message": message,
				"numbers": list.CallList,
				"error":   err.Error(),
			}).Error("Hoo. Couldn't create new Call task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		_, err = connections.RedisClient.Enqueue(callTask)
		if err != nil {
			response.Message = "Call - could not enqueue task"
			log.Log.WithFields(logrus.Fields{
				"message":   message,
				"call_list": list.CallList,
				"task":      callTask,
				"error":     err.Error(),
			}).Error("Hoo. Couldn't enqueue call task!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	data := outputs.Send{
		Message: message,
		SMS:     list.SMSList,
		Call:    list.CallList,
	}
	response.Data = data
	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"message":   message,
		"sms_list":  list.SMSList,
		"call_list": list.CallList,
		"response":  response.Message,
	}).Info("Hoo. Create sending jobs successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

func CheckSeverity(contacts []models.HooQuery, severity int) []models.HooQuery {
	var validContacts []models.HooQuery

	for _, value := range contacts {
		if severity >= value.Severity {
			validContacts = append(validContacts, value)
		}
	}

	return validContacts
}

func GetSendingList(contacts []models.HooQuery, alert_message string, values []string) (*models.SendList, string, []models.Log, error) {
	var contactIDs []int
	var logList []models.Log

	logBody.AlertGroupID = contacts[0].AlertGroupID
	logBody.ProjectAlertGroupID = contacts[0].ProjectAlertGroupID

	for _, value := range contacts {
		contactIDs = append(contactIDs, value.ContactID)
	}

	contactsInfo, err := repositories.GetContactList(contactIDs)
	if err != nil {
		return nil, "", nil, err
	}

	var SMSList []string
	var callList []string
	var list models.SendList

	tmp := make([]interface{}, len(values))
	for i, val := range values {
		tmp[i] = val
	}
	message := fmt.Sprintf(alert_message, tmp...)
	logBody.Message = message

	for _, value := range contacts {

		if value.SMS {
			logRecord := models.Log{
				UserID:              logBody.UserID,
				AlertID:             logBody.AlertID,
				AlertGroupID:        logBody.AlertGroupID,
				ProjectAlertGroupID: logBody.ProjectAlertGroupID,
				Message:             logBody.Message,
				Service:             logBody.Service,
				ProjectID:           logBody.ProjectID,
			}
			logRecord.Channel = constants.SMS
			for _, val := range contactsInfo {
				if val.ID == value.ContactID {
					logRecord.ContactID = value.ContactID
					logRecord.ContactProjectAlertGroupID = val.ID
					logRecord.ContactInfo = val.PhoneNumber
					logList = append(logList, logRecord)
					SMSList = append(SMSList, val.PhoneNumber)
				}
			}
		}
		if value.Call {
			logRecord := models.Log{
				UserID:              logBody.UserID,
				AlertID:             logBody.AlertGroupID,
				AlertGroupID:        logBody.AlertGroupID,
				ProjectAlertGroupID: logBody.ProjectAlertGroupID,
				Message:             logBody.Message,
				Service:             logBody.Service,
				ProjectID:           logBody.ProjectID,
			}
			logRecord.Channel = constants.CALL
			for _, val := range contactsInfo {
				if val.ID == value.ContactID {
					logRecord.ContactID = value.ContactID
					logRecord.ContactProjectAlertGroupID = val.ID
					logRecord.ContactInfo = val.PhoneNumber
					logList = append(logList, logRecord)
					callList = append(callList, val.PhoneNumber)
				}
			}
		}
	}
	list.SMSList = SMSList
	list.CallList = callList
	return &list, message, logList, nil

}
