package controllers

import (
	"regexp"
	"strconv"

	gopkgs "github.com/abr-ooo/go-pkgs"
	"github.com/abr-ooo/hodhod/inputs"
	"github.com/abr-ooo/hodhod/models"
	"github.com/abr-ooo/hodhod/repositoies"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/abr-ooo/hodhod/log"
	"github.com/sirupsen/logrus"
)

// GetUserContact godoc
// @Summary get user contacts
// @Description return contact info
// @ID get_user_contacts
// @Tags Contact
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /contact [get]
func GetUserContacts(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	userID := GetUserID(c)

	contacts, err := repositories.GetContacts(userID)
	if err != nil {
		response.Message = "User Have No Contact!"
		log.Log.WithFields(logrus.Fields{
			"userID":   userID,
			"contacts": contacts,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("Getcontacts. There is no contact for this user!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contacts
	log.Log.WithFields(logrus.Fields{
		"userID":   userID,
		"contacts": contacts,
		"response": response.Message,
	}).Info("GetContacts. Get contacts successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// GetContact godoc
// @Summary get contact
// @Description return contact info
// @ID get_contact
// @Tags Contact
// @Param id path int true "contact_id"
// @Security ApiKeyAuth
// @Success 200 {object} []models.Response
// @Failure 400 json httputil.HTTPError
// @Router /contact/{id} [get]
func GetContact(c *fiber.Ctx) error {
	var response models.Response
	response.Status = "error"
	contactID, err := c.ParamsInt("id")
	if err != nil {
		response.Message = "Parse Contact ID Failed!"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"input":     c.Params("id"),
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("Getcontact. Parse Contact ID to int Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contact, err := repositories.GetContact(contactID)
	if err != nil {
		response.Message = "This Contact does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"contact":   contact,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("Getcontact. This Contact Doesn't Exist In DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, contact.UserID)
	if !access {
		response.Message = "User Doesn't Access To Contact"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      contact.UserID,
			"response":    response.Message,
		}).Info("Getcontact. This user doesn't create this Contact and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contact
	log.Log.WithFields(logrus.Fields{
		"userID":   userID,
		"contact":  contact,
		"response": response.Message,
	}).Info("GetContact. Get contact successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// NewContact godoc
// @Summary new contact
// @Description new contact
// @ID new_contact
// @Tags Contact
// @Param contactBody body inputs.ContactBody true "contact"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /contact [post]
func NewContact(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"

	contactBody := new(inputs.ContactBody)
	err := c.BodyParser(contactBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContact. Parse body to contact body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(contactBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("NewContact. Validate Contact body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var existName, existEmail, existPhoneNumber bool = contactBody.Name != "", contactBody.Email != "", contactBody.PhoneNumber != ""

	if !existName {
		response.Message = "Name Can't empty"
		log.Log.WithFields(logrus.Fields{
			"phoneNumber": contactBody.PhoneNumber,
			"response":    response.Message,
		}).Info("NewContact. Name is empty!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if existPhoneNumber {
		validPhoneNumber, err := Validation("PHONE", contactBody.PhoneNumber)
		if err != nil {
			response.Message = "Check phone number Format Failed"
			log.Log.WithFields(logrus.Fields{
				"phoneNumber": contactBody.PhoneNumber,
				"response":    response.Message,
				"error":       err.Error(),
			}).Error("NewContact. Validation function of phone number have error!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		if !*validPhoneNumber {
			response.Message = "phone number Format Doesn't Valid"
			log.Log.WithFields(logrus.Fields{
				"phoneNumber": contactBody.PhoneNumber,
				"response":    response.Message,
			}).Info("NewContact. This phone number format doesn't valid!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if existEmail {
		validEmail, err := Validation("EMAIL", contactBody.Email)
		if err != nil {
			response.Message = "Check email Format Failed"
			log.Log.WithFields(logrus.Fields{
				"email":    contactBody.Email,
				"response": response.Message,
				"error":    err.Error(),
			}).Error("NewContact. Validation function of email have error!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		if !*validEmail {
			response.Message = "email Format Doesn't Valid"
			log.Log.WithFields(logrus.Fields{
				"email":    contactBody.Email,
				"response": response.Message,
			}).Info("NewContact. This email format doesn't valid!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	userID := GetUserID(c)
	exist := repositories.ExistContact(userID, *contactBody)
	if exist {
		response.Message = "Duplicate Contact"
		log.Log.WithFields(logrus.Fields{
			"name":        contactBody.Name,
			"phoneNumber": contactBody.PhoneNumber,
			"Email":       contactBody.Email,
			"userID":      userID,
			"response":    response.Message,
		}).Info("NewContact. This contact is duplicate. we have one of this!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contact, err := repositories.CreateContact(userID, *contactBody)
	if err != nil {
		response.Message = "Insert Contact In DB Failed"
		log.Log.WithFields(logrus.Fields{
			"contact_model": contact,
			"contact_body":  contactBody,
			"response":      response.Message,
			"error":         err.Error(),
		}).Error("NewContact. Insert Contact In DB Have Error!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	response.Data = contact.ID
	log.Log.WithFields(logrus.Fields{
		"contact":  contact,
		"response": response.Message,
	}).Info("NewContact. Create contact successful :)")
	return c.Status(fiber.StatusOK).JSON(response)

}

// UpdateContact godoc
// @Summary update contact
// @Description update contact
// @ID update_contact
// @Tags Contact
// @Param id path int true "contact_id"
// @Param contactBody body inputs.ContactBody true "contact"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /contact/{id} [put]
func UpdateContact(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	contactID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		response.Message = "Parse Contact ID Failed!"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"input":     c.Params("id"),
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("Getcontact. Parse Contact ID to int Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contactBody := new(inputs.ContactBody)
	err = c.BodyParser(contactBody)
	if err != nil {
		response.Message = "Parse Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("َUpdateContact. Parse body to contact body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	err = validate.Struct(contactBody)
	if err != nil {
		response.Message = "Validate Body Failed"
		log.Log.WithFields(logrus.Fields{
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateContact. Validate Contact body failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var existName, existEmail, existPhoneNumber bool = contactBody.Name != "", contactBody.Email != "", contactBody.PhoneNumber != ""

	if existPhoneNumber {
		validPhoneNumber, err := Validation("PHONE", contactBody.PhoneNumber)
		if err != nil {
			response.Message = "Check phone number Format Failed"
			log.Log.WithFields(logrus.Fields{
				"phoneNumber": contactBody.PhoneNumber,
				"response":    response.Message,
				"error":       err.Error(),
			}).Error("UpdateContact. Validation function of phone number have error!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		if !*validPhoneNumber {
			response.Message = "phone number Format Doesn't Valid"
			log.Log.WithFields(logrus.Fields{
				"phoneNumber": contactBody.PhoneNumber,
				"response":    response.Message,
			}).Info("UpdateContact. This phone number format doesn't valid!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if existEmail {
		validEmail, err := Validation("EMAIL", contactBody.Email)
		if err != nil {
			response.Message = "Check email Format Failed"
			log.Log.WithFields(logrus.Fields{
				"email":    contactBody.Email,
				"response": response.Message,
				"error":    err.Error(),
			}).Error("UpdateContact. Validation function of email have error!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		if !*validEmail {
			response.Message = "email Format Doesn't Valid"
			log.Log.WithFields(logrus.Fields{
				"email":    contactBody.Email,
				"response": response.Message,
			}).Info("UpdateContact. This email format doesn't valid!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	contact, err := repositories.GetContact(contactID)
	if err != nil {
		response.Message = "This Contact does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"contact":   contact,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("Updatecontact. This Contact Doesn't Exist In DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, contact.UserID)
	if !access {
		response.Message = "User Doesn't Access To Contact"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      contact.UserID,
			"response":    response.Message,
		}).Info("Updatecontact. This user doesn't create this Contact and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if existName {
		existContactName := repositories.ExistContact(userID, *contactBody)
		if existContactName {
			response.Message = "Duplicate Contact(Same Name)!"
			log.Log.WithFields(logrus.Fields{
				"name":     contact.Name,
				"userID":   userID,
				"response": response.Message,
			}).Info("UpdateContact. This contact is duplicate. we have one of this name!")
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	err = repositories.UpdateContact(*contact, contactID, *contactBody, existName, existEmail, existPhoneNumber)
	if err != nil {
		response.Message = "Update Contact Failed"
		log.Log.WithFields(logrus.Fields{
			"contact":  contact,
			"response": response.Message,
			"error":    err.Error(),
		}).Error("UpdateContact. update contact in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"contact":  contact,
		"response": response.Message,
	}).Info("UpdateContact. Update contact successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteContact godoc
// @Summary delete contact
// @Description delete contact by contactID
// @ID delete_contact
// @Tags Contact
// @Param id path int true "contact_id"
// @Security ApiKeyAuth
// @Success 200 {object} models.Response
// @Failure 400 json httputil.HTTPError
// @Router /contact/{id} [delete]
func DeleteContact(c *fiber.Ctx) error {

	var response models.Response
	response.Status = "error"
	contactID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		response.Message = "Parse Contact ID Failed!"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"input":     c.Params("id"),
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("Getcontact. Parse Contact ID to int Failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	contact, err := repositories.GetContact(contactID)
	if err != nil {
		response.Message = "This Contact does not exist"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"contact":   contact,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("DeleteContact. This Contact Doesn't Exist In DB!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, access := CheckAccess(c, contact.UserID)
	if !access {
		response.Message = "User Doesn't Access To Contact"
		log.Log.WithFields(logrus.Fields{
			"tokenUserID": userID,
			"userID":      contact.UserID,
			"response":    response.Message,
		}).Info("DeleteContact. This user doesn't create this Contact and doesn't have access!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	err = repositories.DeleteContact(*contact)
	if err != nil {
		response.Message = "Delete Contact Failed"
		log.Log.WithFields(logrus.Fields{
			"contactID": contactID,
			"response":  response.Message,
			"error":     err.Error(),
		}).Error("DeleteContact. delete contact in DB failed!")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	response.Message = "OK!"
	response.Status = "succes"
	log.Log.WithFields(logrus.Fields{
		"contactID": contactID,
		"response":  response.Message,
	}).Info("DeleteContact. Delete contact successful :)")
	return c.Status(fiber.StatusOK).JSON(response)
}

// func VerifyUser(table string, auth models.AuthResponse) (bool, string) {
// 	userID := strconv.Itoa(auth.UserID)

// 	exist := repositories.Exist(table, "user_id", userID)
// 	log.Log.WithFields(logrus.Fields{
// 		"userID": userID,
// 		"exist":  exist,
// 	}).Debug("VerifyUser. success!")
// 	return exist, userID
// }

func Validation(valueType string, value string) (*bool, error) {

	var r *regexp.Regexp
	var err error

	switch valueType {
	case "PHONE":
		r, err = regexp.Compile(`^(\+98|0|98|0098)?([ ]|-|[()]){0,2}9[0-9]([ ]|-|[()]){0,2}(?:[0-9]([ ]|-|[()]){0,2}){8}$`)
	case "EMAIL":
		r, err = regexp.Compile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

	}
	match := r.MatchString(value)

	if err != nil {
		log.Log.WithFields(logrus.Fields{
			"recordType":  valueType,
			"recordValue": value,
			"error":       err.Error(),
		}).Debug("Aux.Validation have error!")
		return nil, err
	}

	log.Log.WithFields(logrus.Fields{
		"recordType":  valueType,
		"recordValue": value,
		"match":       match,
	}).Debug("Aux.Validation finish :))")
	return &match, nil

}

func CheckAccess(c *fiber.Ctx, userID int) (int, bool) {

	tokenID := GetUserID(c)
	return tokenID, tokenID == userID

}

func GetUserID(c *fiber.Ctx) int {

	return int(gopkgs.UID(c))
	

}
