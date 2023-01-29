package repositories

import (
	"errors"
	"strconv"

	"github.com/mhd7966/hodhod/connections"
	"github.com/mhd7966/hodhod/inputs"
	"github.com/mhd7966/hodhod/models"
	"github.com/mhd7966/hodhod/outputs"
	"gorm.io/gorm/clause"
)

func ExistAlertGroup(userID int, name string) bool {
	var alertGroup models.AlertGroup

	if result := connections.DB.Where(models.AlertGroup{
		UserID: userID,
		Name:   name,
	}).First(&alertGroup); result.Error != nil {
		return false
	}
	return true
}

func GetAlertGroup(alertGroupID int) (*models.AlertGroup, error) {
	var alertGroup models.AlertGroup

	if result := connections.DB.First(&alertGroup, alertGroupID); result.Error != nil {
		return nil, result.Error
	}

	return &alertGroup, nil
}

func GetAlertGroupWithDetail(alertGroupID int) (int, *outputs.AlertGroupOutput, error) {
	var alertGroup outputs.AlertGroupOutput
	type Result struct {
		Name      string
		UserID    int
		IsDefault bool
		Message   string
	}
	var result []Result

	_ = connections.DB.Table("alert_groups").Select("alert_groups.user_id, alert_groups.name, alert_groups.is_default, alerts.message").
		Joins("join alert_group_alerts on alert_groups.id = alert_group_alerts.alert_group_id and alert_groups.id = " + strconv.Itoa(alertGroupID)).
		Joins("join alerts on alerts.id = alert_group_alerts.alert_id").
		Scan(&result)

	if len(result) < 1 {
		return 0, nil, errors.New("there isn't any alert group")
	}

	alertGroup.Name = result[0].Name
	alertGroup.IsDefault = result[0].IsDefault
	for _, value := range result {
		alertGroup.Alerts = append(alertGroup.Alerts, value.Message)
	}

	return result[0].UserID, &alertGroup, nil
}

func GetAlertGroups(userID int) (*[]models.AlertGroup, error) {
	var alertGroups []models.AlertGroup

	if result := connections.DB.Find(&alertGroups, models.AlertGroup{UserID: userID}); result.RowsAffected < 1 {
		return nil, result.Error
	}

	return &alertGroups, nil
}

func CreateAlertGroup(userID int, input inputs.AlertGroupBody) (*models.AlertGroup, error) {
	alertGroup := models.AlertGroup{
		UserID:    userID,
		Name:      input.Name,
		IsDefault: input.IsDefault,
	}

	if result := connections.DB.Create(&alertGroup); result.Error != nil {
		return nil, result.Error
	}

	var alert_group_alerts []models.AlertGroupAlerts
	var alert_group_alert models.AlertGroupAlerts
	for _, value := range input.AlertIDs {
		alert_group_alert.AlertGroupID = int(alertGroup.ID)
		alert_group_alert.AlertID = value
		alert_group_alerts = append(alert_group_alerts, alert_group_alert)
	}

	if result := connections.DB.Create(&alert_group_alerts); result.Error != nil {
		return nil, result.Error
	}

	return &alertGroup, nil
}

func UpdateAlertGroup(alertGroupID int, input inputs.AlertGroupUpdate) error {

	if result := connections.DB.Model(&models.AlertGroup{}).Where("id = ?", alertGroupID).Update("name", input.Name); result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateAlertGroupAlerts(alertGroupID int, input inputs.AlertGroupAlerts) error {

	if result := connections.DB.Where("alert_group_id = ?", alertGroupID).Delete(&models.AlertGroupAlerts{}); result.Error != nil {
		return result.Error
	}

	var alert_group_alerts []models.AlertGroupAlerts
	var alert_group_alert models.AlertGroupAlerts
	for _, value := range input.AlertIDs {
		alert_group_alert.AlertGroupID = alertGroupID
		alert_group_alert.AlertID = value
		alert_group_alerts = append(alert_group_alerts, alert_group_alert)
	}

	if result := connections.DB.Create(&alert_group_alerts); result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteAlertGroup(alertGroupID int) error {

	if result := connections.DB.Where("id = ?", alertGroupID).Delete(&models.AlertGroup{}); result.Error != nil {
		return result.Error
	}

	var projects []models.ProjectAlertGroup
	if result := connections.DB.Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).Where("alert_group_id = ?", alertGroupID).Delete(&projects); result.Error != nil {
		return result.Error
	}

	var ids []int
	for _, value := range projects {
		ids = append(ids, int(value.ID))
	}

	connections.DB.Delete(&models.ContactProjectAlertGroup{}, ids)

	return nil
}
