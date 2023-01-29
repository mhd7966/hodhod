package repositories

import (
	"errors"
	"fmt"

	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/models"
	"github.com/abr-ooo/hodhod/outputs"
)

func GetLog(logID int) (*models.Log, error) {
	var log models.Log

	if result := connections.DB.First(&log, logID); result.Error != nil {
		return nil, result.Error
	}

	return &log, nil
}

func GetLogs(userID int) (*[]models.Log, error) {
	var logs []models.Log

	if result := connections.DB.Find(&logs, models.Contact{UserID: userID}); result.RowsAffected < 1 {
		return nil, errors.New("user doesn't have any log")
	}

	return &logs, nil
}

func CreateLog(log models.Log) (*models.Log, error) {

	if result := connections.DB.Create(&log); result.Error != nil {
		return nil, result.Error
	}
	return &log, nil
}

func CreateBulkLog(logs []models.Log) ([]models.Log, error) {

	if result := connections.DB.Create(&logs); result.Error != nil {
		return nil, result.Error
	}
	return logs, nil
}

func UpdateLogStatus(logID int, status int) error {

	if result := connections.DB.Model(&models.Log{}).Where("id = ?", logID).Update("status", status); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetLogInfo(userID int) (*[]outputs.LogInfo, error) {
	var logs []outputs.LogInfo

	if result := connections.
		DB.Raw(`select c."name" contact, ag."name" alert_group, pag."name" project_alert_group, l.message,l.service, 
				l.channel, l.status, l.contact_info, l.created_at, l.updated_at from logs l 
				join contacts c on l.contact_id = c.id 
				join alert_groups ag on l.alert_group_id = ag.id 
				join project_alert_groups pag on l.project_alert_group_id = pag.id 
				WHERE l.user_id = ?`, userID).
		Scan(&logs); result.RowsAffected < 1 {
		return nil, result.Error
	}
	fmt.Println(logs)

	return &logs, nil
}
