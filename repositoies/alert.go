package repositories

import (
	"github.com/mhd7966/hodhod/connections"
	"github.com/mhd7966/hodhod/inputs"
	"github.com/mhd7966/hodhod/models"
	// log "github.com/sirupsen/logrus"
)

func ExistAlert(input inputs.AlertBody) bool {
	var alert models.Alert

	if result := connections.DB.Where(models.Alert{
		Message:  input.Message,
		Severity: input.Severity,
		ParentID: input.ParentID,
	}).First(&alert); result.Error != nil {
		return false
	}
	return true
}

func ParentIsCorrect(parentID int) bool {
	var alert models.Alert

	if result := connections.DB.Where("parent_id = ? ", 0).Where("id = ?", parentID).First(&alert); result.Error != nil {
		return false
	}
	return true

}

func GetAlerts() (*[]models.Alert, error) {
	var alerts []models.Alert

	if result := connections.DB.Find(&alerts); result.RowsAffected < 1 {
		return nil, result.Error
	}

	return &alerts, nil
}

func GetAlert(alertID int) (*models.Alert, error) {
	var alert models.Alert

	if result := connections.DB.First(&alert, alertID); result.Error != nil {
		return nil, result.Error
	}

	return &alert, nil
}

func CreateAlert(input inputs.AlertBody) (*models.Alert, error) {
	alert := models.Alert{
		Message:  input.Message,
		Severity: input.Severity,
		ParentID: input.ParentID,
	}

	if result := connections.DB.Create(&alert); result.Error != nil {
		return nil, result.Error
	}
	return &alert, nil
}

func UpdateAlert(alertID int, alert models.Alert, input inputs.AlertUpdate, message bool, severity bool) error {

	if message {
		alert.Message = input.Message
	}
	if severity {
		alert.Severity = input.Severity
	}

	if result := connections.DB.Save(&alert); result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteAlert(alertID int) error {

	if result := connections.DB.Where("id = ?", alertID).Delete(&models.Alert{}); result.Error != nil {
		return result.Error
	}

	if result := connections.DB.Where("alert_group_id = ?", alertID).Delete(&models.AlertGroupAlerts{}); result.Error != nil {
		return result.Error
	}

	return nil
}
