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

// GetAlertGroups godoc
// @Summary get AlertGroups
// @Description return AlertGroups
// @ID get_user_alertGroups
// @Tags Alert Group
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group [get]
func GetAlertGroups(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	userID := GetUserID(c)

	alertGroups, err := repositories.GetAlertGroups(userID)
	if err != nil {
		response.Message = "This User Have No Alert Group!"
		log.Log.WithFields(logrus.Fields{
			"alertGroup": alertGroups,
			"response":   response.Message,
			"error":      err.Error(),
		}).Error("GetAlertGroups. This User Have No Alert Group in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alertGroups
	log.Log.WithFields(logrus.Fields{
		"alertGroups": alertGroups,
		"response":    response.Message,
	}).Info("GetAlertGroups. Get Alert Groups Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetAlertGroup godoc
// @Summary get AlertGroup
// @Description return AlertGroup
// @ID get_alert_Group
// @Tags Alert Group
// @Param id path int true "alert_grp_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id} [get]
func GetAlertGroup(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("GetAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	dbUserID, alertGroups, err := repositories.GetAlertGroupWithDetail(alertGroupID)
	if err != nil {
		response.Message = "This Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"alertGroups":  alertGroups,
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("GetAlertGroup. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, dbUserID)
	if !access {
		response.Message = "User Doesn't Access To Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      dbUserID,
			"response":    response.Message,
		}).Info("GetAlertGroup. This user doesn't create this Alert Group and doesn't have access to Get!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alertGroups
	log.Log.WithFields(logrus.Fields{
		"userID":     userID,
		"alertGroup": alertGroups,
		"response":   response.Message,
	}).Info("GetAlertGroup. Get Alert Group successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// NewAlertGroup godoc
// @Summary new AlertGroup
// @Description new AlertGroup
// @ID new_AlertGroup
// @Tags Alert Group
// @Param alert-group body inputs.AlertGroupBody true "alert_group"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group [post]
func NewAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	alertGroupBody := new(inputs.AlertGroupBody)
	err := c.BodyParser(alertGroupBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewAlertGroup. Parse body to alertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(alertGroupBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewAlertGroup. Validate AlertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID := GetUserID(c)

	exist := repositories.ExistAlertGroup(userID, alertGroupBody.Name)
	if exist {
		response.Message = "Duplicate AlertGroup Name"
		log.Log.WithFields(logrus.Fields{
			"name":     alertGroupBody.Name,
			"userID":   userID,
			"response": response.Message,
		}).Info("NewAlertGroup. This alert group name is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertGroup, err := repositories.CreateAlertGroup(userID, *alertGroupBody)
	if err != nil {
		response.Message = "Insert AlertGroup In DB Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroup_body":  alertGroupBody,
			"alertGroup_model": alertGroup,
			"response":         response.Message,
			"error":            err.Error(),
		}).Error("NewAlertGroup. Insert AlertGroup In DB Have Error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = alertGroup.ID
	log.Log.WithFields(logrus.Fields{
		"alertGroup": alertGroup,
		"response":   response.Message,
	}).Info("NewAlertGroup. Create AlertGroup successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// UpdateAlertGroup godoc
// @Summary update AlertGroup
// @Description update AlertGroup
// @ID update_AlertGroup
// @Tags Alert Group
// @Param id path int true "alert_group_id"
// @Param alert_group body inputs.AlertGroupUpdate true "alert_group"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id} [put]
func UpdateAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("UpdateAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertGroupBody := new(inputs.AlertGroupUpdate)
	err = c.BodyParser(alertGroupBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroup": alertGroupBody,
			"response":   response.Message,
			"error":      err.Error(),
		}).Error("َUpdateAlertGroup. Parse body to alertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(alertGroupBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateAlertGroup. Validate AlertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertGroup, err := repositories.GetAlertGroup(alertGroupID)
	if err != nil {
		response.Message = "This Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"alertGroup":   alertGroup,
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("َUpdateAlertGroup. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, alertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      alertGroup.UserID,
			"response":    response.Message,
		}).Info("َUpdateAlertGroup. This user doesn't create this Alert Group and doesn't have access to update!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	existAlertGroupName := repositories.ExistAlertGroup(userID, alertGroupBody.Name)
	if existAlertGroupName {
		response.Message = "Duplicate AlertGroup(Same Name)!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupName": alertGroupBody.Name,
			"userID":         userID,
			"response":       response.Message,
		}).Info("UpdateAlertGroup. This alert group is duplicate. we have one of this name!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.UpdateAlertGroup(alertGroupID, *alertGroupBody)
	if err != nil {
		response.Message = "Update AlertGroup Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroup": alertGroupBody,
			"response":   response.Message,
			"error":      err.Error(),
		}).Error("UpdateAlertGroup. update alertGroup in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"alertGroup": alertGroupBody,
		"response":   response.Message,
	}).Info("UpdateAlertGroup. Update alertGroup successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteAlertGroup godoc
// @Summary delete alertGroup
// @Description delete alertGroup by alertGroupID
// @ID delete_alert_group
// @Tags Alert Group
// @Param id path int true "alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id} [delete]
func DeleteAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("DeleteAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertGroup, err := repositories.GetAlertGroup(alertGroupID)
	if err != nil {
		response.Message = "This Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"alertGroup":   alertGroup,
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("DeleteAlertGroup. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, alertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      alertGroup.UserID,
			"response":    response.Message,
		}).Info("DeleteAlertGroup. This user doesn't create this Alert Group and doesn't have access to delete!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.DeleteAlertGroup(alertGroupID)
	if err != nil {
		response.Message = "Delete AlertGroup Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("DeleteAlertGroup. delete AlertGroup in DB Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"alertGroupID": alertGroupID,
		"response":     response.Message,
	}).Info("DeleteAlertGroup. Delete AlertGroup Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateAlertGroupAlerts godoc
// @Summary update AlertGroupAlerts
// @Description update AlertGroupAlerts
// @ID update_AlertGroupAlerts
// @Tags Alert Group
// @Param id path int true "alert_group_id"
// @Param alerts body inputs.AlertGroupAlerts true "alerts"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/alert [put]
func UpdateAlertGroupAlerts(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("UpdateAlertGroupAlerts. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertsBody := new(inputs.AlertGroupAlerts)
	err = c.BodyParser(alertsBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroup": alertsBody,
			"response":   response.Message,
			"error":      err.Error(),
		}).Error("َUpdateAlertGroupAlerts. Parse body to alertGroupAlerts body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	alertGroup, err := repositories.GetAlertGroup(alertGroupID)
	if err != nil {
		response.Message = "This Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"alertGroup":   alertGroup,
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("َUpdateAlertGroupAlerts. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, alertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      alertGroup.UserID,
			"response":    response.Message,
		}).Info("َUpdateAlertGroupAlerts. This user doesn't create this Alert Group and doesn't have access to update!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.UpdateAlertGroupAlerts(alertGroupID, *alertsBody)
	if err != nil {
		response.Message = "Update AlertGroup Alerts Failed"
		log.Log.WithFields(logrus.Fields{
			"alertGroup": alertsBody,
			"response":   response.Message,
			"error":      err.Error(),
		}).Error("UpdateAlertGroupAlerts. update alertGroup alerts in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"response": response.Message,
	}).Info("UpdateAlertGroupAlerts. Update alertGroup alert successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}
