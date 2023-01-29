package controllers

import (
	"github.com/abr-ooo/hodhod/inputs"
	"github.com/abr-ooo/hodhod/models"
	"github.com/abr-ooo/hodhod/repositoies"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/abr-ooo/hodhod/log"

)

// GetAlerts godoc
// @Summary get alerts
// @Description return alerts
// @ID get_alerts
// @Tags Alert
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert [get]
func GetAlerts(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	alerts, err := repositories.GetAlerts()
	if err != nil {
		response.Message = "Get alerts Failed!"
		log.Log.WithFields(logrus.Fields{
			"alerts":   alerts,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetAlerts. Get alerts from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alerts
	log.Log.WithFields(logrus.Fields{
		"alerts":   alerts,
		"response": response.Message,
	}).Info("GetAlerts. Get alerts successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetAlert godoc
// @Summary get alert
// @Description return alert
// @ID get_alert
// @Tags Alert
// @Param id path int true "alert_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert/{id} [get]
func GetAlert(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	alertID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetAlert. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alert, err := repositories.GetAlert(alertID)
	if err != nil {
		response.Message = "Get alert Failed!"
		log.Log.WithFields(logrus.Fields{
			"alert":    alert,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("GetAlert. Get alert from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alert
	log.Log.WithFields(logrus.Fields{
		"alerts":   alert,
		"response": response.Message,
	}).Info("GetAlert. Get alert successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// NewAlert godoc
// @Summary new alert
// @Description create alert
// @ID new_alert
// @Tags Alert
// @Param alertBody body inputs.AlertBody true "alert"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert [post]
func NewAlert(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	alertBody := new(inputs.AlertBody)
	err := c.BodyParser(alertBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewAlert. Parse body to alert body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(alertBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewAlert. Validate Alert body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var existSeverity bool = alertBody.Severity != 0
	if alertBody.ParentID != 0 && !existSeverity {
		response.Message = "Severity Is Empty!"
		log.Log.WithFields(logrus.Fields{
			"message":  alertBody.Message,
			"severity": alertBody.Severity,
			"response": response.Message,
		}).Info("NewAlert. Child ALert Must Have Severity!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if alertBody.ParentID == 0 && existSeverity {
		response.Message = "Severity Is Fill!"
		log.Log.WithFields(logrus.Fields{
			"message":  alertBody.Message,
			"severity": alertBody.Severity,
			"response": response.Message,
		}).Info("NewAlert. Parent ALert Must NOT Have Severity!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	exist := repositories.ExistAlert(*alertBody)
	if exist {
		response.Message = "Duplicate Alert"
		log.Log.WithFields(logrus.Fields{
			"message":  alertBody.Message,
			"severity": alertBody.Severity,
			"response": response.Message,
		}).Info("NewAlert. This alert is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if alertBody.ParentID != 0 {
		ok := repositories.ParentIsCorrect(alertBody.ParentID)
		if !ok {
			response.Message = "Alert ParentID Is Wrong"
			log.Log.WithFields(logrus.Fields{
				"PID":      alertBody.ParentID,
				"response": response.Message,
			}).Info("NewAlert. This parent id this isn't parent!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	alert, err := repositories.CreateAlert(*alertBody)
	if err != nil {
		response.Message = "Insert Alert In DB Failed"
		log.Log.WithFields(logrus.Fields{
			"alert_model": alert,
			"alert_body":  alertBody,
			"response":    response.Message,
			"error":       err.Error(),
		}).Error("NewAlert. Insert Alert In DB Have Error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alert.ID
	log.Log.WithFields(logrus.Fields{
		"alert":    alert,
		"response": response.Message,
	}).Info("NewAlert. Create Alert successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// UpdateAlert godoc
// @Summary update alert
// @Description update alert
// @ID update_alert
// @Tags Alert
// @Param id path int true "alert_id"
// @Param alertBody body inputs.AlertUpdate true "alert"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert/{id} [put]
func UpdateAlert(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	alertID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateAlert. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertBody := new(inputs.AlertUpdate)
	err = c.BodyParser(alertBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateAlert. Parse body to alert body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(alertBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateAlertGroup. Validate Alert body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var existMessage, existSeverity bool = alertBody.Message != "", alertBody.Severity != 0

	if !existMessage && !existSeverity {
		response.Message = "Validate Body Failed. At least one field must be entered!"
		log.Log.WithFields(logrus.Fields{
			"existMessage":  existMessage,
			"existSeverity": existSeverity,
			"response":      response.Message,
			"error":         err.Error(),
		}).Error("UpdateAlertGroup. Validate Alert body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alert, err := repositories.GetAlert(alertID)
	if err != nil {
		response.Message = "Get alert Failed!"
		log.Log.WithFields(logrus.Fields{
			"alert":    alert,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateAlert. Get alert from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.UpdateAlert(alertID, *alert, *alertBody, existMessage, existSeverity)
	if err != nil {
		response.Message = "Update Contact Failed"
		log.Log.WithFields(logrus.Fields{
			"alert":          alertBody,
			"exist_message":  existMessage,
			"exist_severity": existSeverity,
			"alertID":        alertID,
			"response":       response.Message,
			"error":          err.Error(),
		}).Error("UpdateAlert. update contact in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alert
	log.Log.WithFields(logrus.Fields{
		"alert":    alert,
		"response": response.Message,
	}).Info("UpdateAlert. Update alert successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// DeleteAlert godoc
// @Summary delete alert
// @Description delete alert
// @ID delete_alert
// @Tags Alert
// @Param id path int true "alert_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert/{id} [delete]
func DeleteAlert(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"

	alertID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("DeleteAlert. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.DeleteAlert(alertID)
	if err != nil {
		response.Message = "delete alert Failed!"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("DeleteAlert. delete alert from DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"response": response.Message,
	}).Info("DeleteAlert. Delete alert successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}
