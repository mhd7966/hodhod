package repositories

import (
	"github.com/abr-ooo/hodhod/connections"
	"github.com/abr-ooo/hodhod/models"
)

func GetHooContact(userID int, alertID int, projectID int, service string) (*[]models.HooQuery, error) {
	var contacts []models.HooQuery

	if result := connections.
		DB.Raw(`select cpag.*, pag.alert_group_id from alert_groups ag join alert_group_alerts aga 
			on ag.id = aga.alert_group_id and ag.user_id = ? and aga.alert_id = ?
			join project_alert_groups pag 
			on  pag.alert_group_id = aga.alert_group_id 
			and pag.project_id = ?  
			and pag.service = ? 
			join contact_project_alert_groups cpag 
			on cpag.project_alert_group_id = pag.id`,
		userID, alertID, projectID, service).
		Scan(&contacts); result.RowsAffected < 1 {
		return nil, result.Error
	}

	return &contacts, nil
}
