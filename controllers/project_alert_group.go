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

// GetProjectsAlertGroup godoc
// @Summary get ProjectsAlertGroup
// @Description return ProjectsAlertGroup
// @ID get_projects_of_AlertGroup
// @Tags Project Alert Group
// @Param id path string true "alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project [get]
func GetProjectsAlertGroup(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID": alertGroupID,
			"response":     response.Message,
			"error":        err.Error(),
		}).Error("GetProjectsAlertGroup. Convert param string to int failed!")
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
		}).Error("GetProjectsAlertGroup. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if !alertGroup.IsDefault {
		userID, access := CheckAccess(c, alertGroup.UserID)
		if !access {
			response.Message = "User Doesn't Access To Alert Group"
			log.Log.WithFields(logrus.Fields{
				"tokenUserID": userID,
				"userID":      alertGroup.UserID,
				"response":    response.Message,
			}).Info("GetProjectsAlertGroup. This user doesn't create this Alert Group and doesn't have access to delete!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	projectsAlertGroup, err := repositories.GetProjectsAlertGroup(alertGroupID)
	if err != nil {
		response.Message = "There Is No Project Of AlertGroup"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID":       alertGroupID,
			"projectsAlertGroup": projectsAlertGroup,
			"response":           response.Message,
			"error":              err.Error(),
		}).Error("GetProjectsAlertGroup. There is no project of this alert group!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if projectsAlertGroup == nil {
		response.Message = "There Is No Project Of AlertGroup"
		log.Log.WithFields(logrus.Fields{
			"alertGroupID":       alertGroupID,
			"projectsAlertGroup": projectsAlertGroup,
			"response":           response.Message,
		}).Error("GetProjectsAlertGroup. There is no project of this alert group!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var a []models.ProjectAlertGroup = *projectsAlertGroup
	if len(*projectsAlertGroup) != 0 && alertGroup.IsDefault {
		userID, access := CheckAccess(c, a[0].UserID)
		if !access {
			response.Message = "User Doesn't Access To Project Alert Group"
			log.Log.WithFields(logrus.Fields{
				"tokenUserID": userID,
				"userID":      alertGroup.UserID,
				"response":    response.Message,
			}).Info("GetProjectsAlertGroup. This user doesn't create this Project Alert Group and doesn't have access to delete!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = projectsAlertGroup
	log.Log.WithFields(logrus.Fields{
		"projectsAlertGroups": projectsAlertGroup,
		"response":            response.Message,
	}).Info("GetProjectsAlertGroups. Get Projects Alert Group Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetProjectAlertGroup godoc
// @Summary get ProjectAlertGroup
// @Description return ProjectAlertGroup
// @ID get_ProjectAlertGroup
// @Tags Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid} [get]
func GetProjectAlertGroup(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("GetProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Proejct Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"proejctAlertGroupID": projectAlertGroupID,
			"projectAlertGroup":   projectAlertGroup,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("GetProjectAlertGroup. This Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"project":     projectAlertGroup,
			"response":    response.Message,
		}).Info("GetProjectAlertGroup. This user doesn't create this Project Alert Group and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = projectAlertGroup
	log.Log.WithFields(logrus.Fields{
		"userID":            projectAlertGroup.UserID,
		"projectAlertGroup": projectAlertGroup,
		"response":          response.Message,
	}).Info("GetProjectAlertGroup. Get Project Alert Group successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// NewProjectAlertGroup godoc
// @Summary new ProjectAlertGroup
// @Description new ProjectAlertGroup
// @ID new_ProjectAlertGroup
// @Tags Project Alert Group
// @Param id path int true "alert_group_id"
// @Param project body inputs.ProjectAlertGroupBody true "project_alert_group"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project [post]
func NewProjectAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	alertGroupID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": alertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("NewProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroupBody := new(inputs.ProjectAlertGroupBody)
	err = c.BodyParser(projectAlertGroupBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewProjectAlertGroup. Parse body to projectAlertGroup model failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(projectAlertGroupBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewProjectAlertGroup. Validate ProjectAlertGroup body failed!")
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
		}).Error("NewProjectAlertGroup. This Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var userID int = 0
	if !alertGroup.IsDefault {
		tokenUserID, access := CheckAccess(c, alertGroup.UserID)
		if !access {
			response.Message = "User Doesn't Access To Alert Group"
			log.Log.WithFields(logrus.Fields{
				"tokenUserID": userID,
				"userID":      alertGroup.UserID,
				"response":    response.Message,
			}).Info("NewProjectAlertGroup. This user doesn't create this Alert Group and doesn't have access to create project!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		temp := &tokenUserID
		userID = *temp
	}

	existProject := repositories.ExistProjectAlertGroup(*projectAlertGroupBody)
	if existProject {
		response.Message = "Duplicate Project "
		log.Log.WithFields(logrus.Fields{
			"project_id": projectAlertGroupBody.ProjectID,
			"service":    projectAlertGroupBody.Service,
			"response":   response.Message,
		}).Info("NewProjectAlertGroup. This project is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	existName := repositories.ExistNameProjectAlertGroup(projectAlertGroupBody.Name, alertGroupID)
	if existName {
		response.Message = "Duplicate Name "
		log.Log.WithFields(logrus.Fields{
			"name":         projectAlertGroupBody.Name,
			"alertGroupID": alertGroupID,
			"response":     response.Message,
		}).Info("NewProjectAlertGroup. This name is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.CreateProjectAlertGroup(userID, alertGroupID, *projectAlertGroupBody)
	if err != nil {
		response.Message = "Insert ProjectAlertGroup In DB Failed"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupBody": projectAlertGroupBody,
			"projectAlertGroup":     projectAlertGroup,
			"userID":                userID,
			"response":              response.Message,
			"error":                 err.Error(),
		}).Error("NewProjectAlertGroup. Insert ProjectAlertGroup In DB Have Error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = projectAlertGroup.ID
	log.Log.WithFields(logrus.Fields{
		"projectAlertGroup": projectAlertGroup,
		"response":          response.Message,
	}).Info("NewProjectAlertGroup. Create ProjectAlertGroup successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// DeleteProjectAlertGroup godoc
// @Summary delete ProjectAlertGroup
// @Description delete ProjectAlertGroup by ProjectAlertGroupID
// @ID delete_project_alert_group
// @Tags Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid} [delete]
func DeleteProjectAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("DeleteProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("َDeleteProjectAlertGroup. This Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"response":    response.Message,
		}).Info("DeleteProjectAlertGroup. This user doesn't create this Alert Group and doesn't have access to delete!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.DeleteProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "Delete ProjectAlertGroup Failed"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("DeleteProjectAlertGroup. Delete ProjectAlertGroup in DB Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"projectAlertGroupID": projectAlertGroupID,
		"response":            response.Message,
	}).Info("DeleteProjectAlertGroup. Delete ProjectAlertGroup Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}
