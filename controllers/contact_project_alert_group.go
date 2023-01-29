package controllers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mhd7966/hodhod/inputs"
	"github.com/mhd7966/hodhod/log"
	"github.com/mhd7966/hodhod/models"
	repositories "github.com/mhd7966/hodhod/repositoies"
	"github.com/sirupsen/logrus"
)

// GetContactProjectsAlertGroup godoc
// @Summary get ContactProjectsAlertGroup
// @Description return ContactProjectsAlertGroup
// @ID get_contacts_of_ProjectAlertGroup
// @Tags Contact Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid}/contact [get]
func GetContactsProjectAlertGroup(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("GetContactsProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("َUpdateContactProjectAlertGroup. This Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"response":    response.Message,
		}).Info("GetContactsProjectsAlertGroup. This user doesn't create this Alert Group and doesn't have access to delete!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactsProjectAlertGroup, err := repositories.GetContactsProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Project Alert Group does not have contact"
		log.Log.WithFields(logrus.Fields{
			"contactsProjectAlertGroup": contactsProjectAlertGroup,
			"projectAlertGroupID":       projectAlertGroupID,
			"response":                  response.Message,
			"error":                     err.Error(),
		}).Error("GetContactsProjectsAlertGroup. This Project Alert Group doesn't have contact in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contactsProjectAlertGroup
	log.Log.WithFields(logrus.Fields{
		"contactsProjectsAlertGroups": contactsProjectAlertGroup,
		"response":                    response.Message,
	}).Info("GetContactsProjectsAlertGroup. Get Contacts projects Alert Group Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetContactprojectAlertGroup godoc
// @Summary get ContactprojectAlertGroup
// @Description return ContactprojectAlertGroup
// @ID get_ContactprojectAlertGroup
// @Tags Contact Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Param cpid path int true "contact_project_alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid}/contact/{cpid} [get]
func GetContactProjectAlertGroup(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	contactProjectAlertGroupID, err := c.ParamsInt("cpid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": contactProjectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("GetContactProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": projectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("GetContactProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactProjectAlertGroup, err := repositories.GetContactProjectAlertGroup(contactProjectAlertGroupID)
	if err != nil {
		response.Message = "This Contact Proejct Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactProejctAlertGroupID": contactProjectAlertGroupID,
			"contactProjectAlertGroup":   contactProjectAlertGroup,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("GetContactProjectAlertGroup. This Contact Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if projectAlertGroupID != contactProjectAlertGroup.ProjectAlertGroupID {
		response.Message = "Information is Wrong"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": contactProjectAlertGroupID,
			"projectAlertGroupID":        projectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("GetContactProjectAlertGroup. Project ALert Group IDs is not same!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This  Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": projectAlertGroup,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("َGetContactProjectAlertGroup. This  Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Contact Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"response":    response.Message,
		}).Info("GetContactProjectAlertGroup. This user doesn't create this Contact Alert Group and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contactProjectAlertGroup
	log.Log.WithFields(logrus.Fields{
		"userID":                   projectAlertGroup.UserID,
		"contactProjectAlertGroup": contactProjectAlertGroup,
		"response":                 response.Message,
	}).Info("GetContactProjectAlertGroup. Get Contact Project Alert Group successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// NewContactProjectALertGroup godoc
// @Summary new ContactProjectALertGroup
// @Description new ContactProjectALertGroup
// @ID new_ContactProjectALertGroup
// @Tags Contact Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Param contact body inputs.NewContactProjectAlertGroup true "contact_project_alert_group"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid}/contact [post]
func NewContactProjectAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("NewContactProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactProjectAlertGroupBody := new(inputs.NewContactProjectAlertGroup)
	err = c.BodyParser(contactProjectAlertGroupBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContactProjectAlertGroup. Parse body to NewcontactProjectAlertGroup failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactUserID, ok, err := CheckOptions(*contactProjectAlertGroupBody)
	if !ok || err != nil {
		response.Message = "Check Options Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContactProjectAlertGroup. Check truth of options failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	existContact := repositories.ExistContactProjectAlertGroup(contactProjectAlertGroupBody.ContactID, projectAlertGroupID)
	if existContact {
		response.Message = "Duplicate Contact "
		log.Log.WithFields(logrus.Fields{
			"contactID":              contactProjectAlertGroupBody.ContactID,
			"project_alert_group_id": projectAlertGroupID,
			"response":               response.Message,
		}).Info("NewContactProjectAlertGroup. This contact is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Contact Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": projectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("َUpdateContactProjectAlertGroup. This Contact Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"response":    response.Message,
		}).Info("NewContactProjectAlertGroup. This user doesn't create this Project Alert Group and doesn't have access to create contact!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if tokenUserID != contactUserID {
		response.Message = "User Doesn't Access To Contact"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID":   tokenUserID,
			"contactUserID": contactUserID,
			"response":      response.Message,
		}).Info("NewContactProjectAlertGroup. This user doesn't create this Contact and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactProjectAlertGroup, err := repositories.CreateContactProjectAlertGroup(projectAlertGroupID, *contactProjectAlertGroupBody)
	if err != nil {
		response.Message = "Insert ContactProjectAlertGroup In DB Failed"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupBody": contactProjectAlertGroupBody,
			"contactProjectAlertGroup":     contactProjectAlertGroup,
			"response":                     response.Message,
			"error":                        err.Error(),
		}).Error("NewContactProjectAlertGroup. Insert ContactProjectAlertGroup In DB Have Error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contactProjectAlertGroup.ID
	log.Log.WithFields(logrus.Fields{
		"projectAlertGroup": contactProjectAlertGroup,
		"response":          response.Message,
	}).Info("NewContactProjectAlertGroup. Create contactProjectAlertGroup successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// UpdateContactProjectAlertGroup godoc
// @Summary update ContactProjectAlertGroup
// @Description update ContactProjectAlertGroup
// @ID update_ContactProjectAlertGroup
// @Tags Contact Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Param cpid path int true "contact_project_alert_group_id"
// @Param contact body inputs.UpdateContactProjectAlertGroup true "contact_project_alert_group"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid}/contact/{cpid} [put]
func UpdateContactProjectAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	contactProjectAlertGroupID, err := c.ParamsInt("cpid")
	if err != nil {
		response.Message = "Convert cpid string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"contactrojectAlertGroupID": contactProjectAlertGroupID,
			"response":                  response.Message,
			"error":                     err.Error(),
		}).Error("UpdateContactProjectAlertGroup. Convert cpid string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	projectAlertGroupID, err := c.ParamsInt("pid")
	if err != nil {
		response.Message = "Convert oid string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"contactrojectAlertGroupID": contactProjectAlertGroupID,
			"response":                  response.Message,
			"error":                     err.Error(),
		}).Error("UpdateContactProjectAlertGroup. Convert oid string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactProjectAlertGroupBody := new(inputs.UpdateContactProjectAlertGroup)
	err = c.BodyParser(contactProjectAlertGroupBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroup": contactProjectAlertGroupBody,
			"response":                 response.Message,
			"error":                    err.Error(),
		}).Error("َUpdateContactProjectAlertGroup. Parse body to updateContactProjectAlertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(contactProjectAlertGroupBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContactProjectAlertGroup. Validate ContactProjectAlertGroup body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactProjectAlertGroup, err := repositories.GetContactProjectAlertGroup(contactProjectAlertGroupID)
	if err != nil {
		response.Message = "This Contact Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": contactProjectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("َUpdateContactProjectAlertGroup. This Contact Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactID := contactProjectAlertGroup.ContactID
	checkOptionsBody := inputs.NewContactProjectAlertGroup{
		ContactID: contactID,
		Severity:  contactProjectAlertGroupBody.Severity,
		Call:      contactProjectAlertGroupBody.Call,
		SMS:       contactProjectAlertGroupBody.SMS,
		Email:     contactProjectAlertGroupBody.Email,
		Webhook:   contactProjectAlertGroupBody.Webhook,
	}

	contactUserID, ok, err := CheckOptions(checkOptionsBody)
	if !ok || err != nil {
		response.Message = "Check Options Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContactProjectAlertGroup. Check truth of options failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// existContact := repositories.ExistContactProjectAlertGroup(contactID, projectAlertGroupID)
	// if existContact {
	// 	response.Message = "Duplicate Contact "
	// 	log.Log.WithFields(logrus.Fields{
	// 		"contactID":              contactID,
	// 		"project_alert_group_id": projectAlertGroupID,
	// 		"response":               response.Message,
	// 	}).Info("UpdateContactProjectAlertGroup. This contact is duplicate. we have one of this!")
	// 	return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	projectAlertGroup, err := repositories.GetProjectAlertGroup(projectAlertGroupID)
	if err != nil {
		response.Message = "This Contact Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": contactProjectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("َUpdateContactProjectAlertGroup. This Contact Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, projectAlertGroup.UserID)
	if !access {
		response.Message = "User Doesn't Access To Contact Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      projectAlertGroup.UserID,
			"response":    response.Message,
		}).Info("َUpdateContactProjectAlertGroup. This user doesn't create this Contact Alert Group and doesn't have access to update!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if tokenUserID != contactUserID {
		response.Message = "User Doesn't Access To Contact"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID":   tokenUserID,
			"contactUserID": contactUserID,
			"response":      response.Message,
		}).Info("UpdateContactProjectAlertGroup. This user doesn't create this Contact and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.UpdateContactProjectAlertGroup(contactProjectAlertGroup, *contactProjectAlertGroupBody)
	if err != nil {
		response.Message = "Update ContactProjectAlertGroup Failed"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroup": contactProjectAlertGroupBody,
			"response":                 response.Message,
			"error":                    err.Error(),
		}).Error("UpdateContactProjectAlertGroup. update ContactProjectAlertGroup in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"contactProjectAlertGroup": contactProjectAlertGroupBody,
		"response":                 response.Message,
	}).Info("UpdateContactProjectAlertGroup. Update ContactProjectAlertGroup successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteContactProjectAlertGroup godoc
// @Summary delete ContactProjectAlertGroup
// @Description delete ContactProjectAlertGroup by ContactProjectAlertGroupID
// @ID delete_contactProjectAlertGroup
// @Tags Contact Project Alert Group
// @Param id path int true "alert_group_id"
// @Param pid path int true "project_alert_group_id"
// @Param cpid path int true "contact_project_alert_group_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /alert-group/{id}/project/{pid}/contact/{cpid} [delete]
func DeleteContactProjectAlertGroup(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	contactProjectAlertGroupID, err := c.ParamsInt("cpid")
	if err != nil {
		response.Message = "Convert param string to int failed!"
		log.Log.WithFields(logrus.Fields{
			"contactrojectAlertGroupID": contactProjectAlertGroupID,
			"response":                  response.Message,
			"error":                     err.Error(),
		}).Error("DeleteContactProjectAlertGroup. Convert param string to int failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, err := repositories.GetContactUserID(contactProjectAlertGroupID)
	if err != nil {
		response.Message = "This Contact Project Alert Group does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactProjectAlertGroupID": contactProjectAlertGroupID,
			"response":                   response.Message,
			"error":                      err.Error(),
		}).Error("َDeleteContactProjectAlertGroup. This Contact Project Alert Group doesn't exist in DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	tokenUserID, access := CheckAccess(c, *userID)
	if !access {
		response.Message = "User Doesn't Access To Contact Project Alert Group"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": tokenUserID,
			"userID":      *userID,
			"response":    response.Message,
		}).Info("DeleteContactProjectAlertGroup. This user doesn't create this Contact Project Alert Group and doesn't have access to delete!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.DeleteContactProjectAlertGroup(contactProjectAlertGroupID)
	if err != nil {
		response.Message = "Delete ContactProjectAlertGroup Failed"
		log.Log.WithFields(logrus.Fields{
			"projectAlertGroupID": contactProjectAlertGroupID,
			"response":            response.Message,
			"error":               err.Error(),
		}).Error("DeleteContactProjectAlertGroup. Delete ContactProjectAlertGroup in DB Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"projectAlertGroupID": contactProjectAlertGroupID,
		"response":            response.Message,
	}).Info("DeleteContactProjectAlertGroup. Delete ContactProjectAlertGroup Successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

func CheckOptions(input inputs.NewContactProjectAlertGroup) (int, bool, error) {

	contact, err := repositories.GetContact(input.ContactID)
	if err != nil {
		return 0, false, err
	}

	if input.Email {
		if contact.Email == "" {
			return 0, false, errors.New("contact doesn't have email address")
		}
	}

	if input.Call || input.SMS {
		if contact.PhoneNumber == "" {
			return 0, false, errors.New("contact doesn't have phone number")
		}
	}

	return contact.UserID, true, nil
}
