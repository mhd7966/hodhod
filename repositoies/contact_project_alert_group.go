package repositories

import (
	"strconv"

	"github.com/mhd7966/hodhod/connections"
	"github.com/mhd7966/hodhod/inputs"
	"github.com/mhd7966/hodhod/models"
	// log "github.com/sirupsen/logrus"
)

func ExistContactProjectAlertGroup(contactID int, projectAlertGroupID int) bool {
	var contactProjectAlertGroup models.ContactProjectAlertGroup

	if result := connections.DB.Where(models.ContactProjectAlertGroup{
		ProjectAlertGroupID: projectAlertGroupID,
		ContactID:           contactID,
	}).First(&contactProjectAlertGroup); result.Error != nil {
		return false
	}
	return true
}

func GetContactsProjectAlertGroup(projectAlertGroupID int) (*[]models.ContactProjectAlertGroup, error) {
	var contactsProjectAlertGroup []models.ContactProjectAlertGroup

	if result := connections.DB.Find(&contactsProjectAlertGroup, models.ContactProjectAlertGroup{ProjectAlertGroupID: projectAlertGroupID}); result.RowsAffected < 1 {
		return nil, result.Error
	}

	return &contactsProjectAlertGroup, nil
}

func GetContactProjectAlertGroup(contactProjectAlertGroupID int) (*models.ContactProjectAlertGroup, error) {
	var contactProjectAlertGroup models.ContactProjectAlertGroup

	if result := connections.DB.First(&contactProjectAlertGroup, contactProjectAlertGroupID); result.Error != nil {
		return nil, result.Error
	}

	return &contactProjectAlertGroup, nil
}

func CreateContactProjectAlertGroup(projectAlertGroupID int, input inputs.NewContactProjectAlertGroup) (*models.ContactProjectAlertGroup, error) {

	contactProjectAlertGroup := models.ContactProjectAlertGroup{
		ProjectAlertGroupID: projectAlertGroupID,
		ContactID:           input.ContactID,
		Severity:            input.Severity,
		Call:                input.Call,
		SMS:                 input.SMS,
		Email:               input.Email,
		Webhook:             input.Webhook,
	}

	if result := connections.DB.Create(&contactProjectAlertGroup); result.Error != nil {
		return nil, result.Error
	}

	return &contactProjectAlertGroup, nil
}

func UpdateContactProjectAlertGroup(contactProjectlertGroup *models.ContactProjectAlertGroup, input inputs.UpdateContactProjectAlertGroup) error {

	contactProjectAlertGroup := models.ContactProjectAlertGroup{
		ID:       contactProjectlertGroup.ID,
		Severity: input.Severity,
		Call:     input.Call,
		SMS:      input.SMS,
		Email:    input.Email,
		Webhook:  input.Webhook,
	}

	if result := connections.DB.Model(contactProjectAlertGroup).Select("Call", "SMS", "Email", "Webhook", "Severity").
		Updates(&contactProjectAlertGroup); result.Error != nil {
		return result.Error
	}

	// if result := connections.DB.Save(&contactProjectAlertGroup); result.Error != nil {
	// 	return result.Error
	// }
	return nil
}

func DeleteContactProjectAlertGroup(contactProjectAlertGroupID int) error {

	if result := connections.DB.Where("id = ?", contactProjectAlertGroupID).Delete(&models.ContactProjectAlertGroup{}); result.Error != nil {
		return result.Error
	}

	return nil
}

func GetContactUserID(contactProjectAlertGroupID int) (*int, error) {

	var userID int
	query := "select t1.user_id from alert_groups t1 inner join project_alert_groups t2 on t1.id = t2.alert_group_id inner join contact_project_alert_groups t3 on t2.id = t3.project_alert_group_id where t3.id=" + strconv.Itoa(contactProjectAlertGroupID)

	if result := connections.DB.Raw(query).Scan(&userID); result.Error != nil {
		return nil, result.Error
	}

	return &userID, nil
}
