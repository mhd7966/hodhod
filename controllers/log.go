package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/log"
	"github.com/mhd7966/hodhod/models"
	repositories "github.com/mhd7966/hodhod/repositoies"
	"github.com/sirupsen/logrus"
)

// Log godoc
// @Summary log
// @Description log
// @ID log
// @Tags Log
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /log [get]
func Log(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	userID := GetUserID(c)

	logs, err := repositories.GetLogInfo(userID)
	if err != nil {
		response.Message = "Get Contact Info Failed"
		log.Log.WithFields(logrus.Fields{
			"user_id":  userID,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("Hoo. There is no contact for this request!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Data = logs
	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"logs":     logs,
		"response": response.Message,
	}).Info("Hoo. Create sending jobs successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}
