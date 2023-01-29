package models

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	ID                         int    `json:"id"`
	UserID                     int    `json:"user_id"`
	ContactID                  int    `json:"contact_id"`
	AlertID                    int    `json:"alert_id"`
	AlertGroupID               int    `json:"alert_group_id"`
	ProjectAlertGroupID        int    `json:"project_alert_group_id"`
	ContactProjectAlertGroupID int    `json:"contact_project_alert_group_id"`
	Message                    string `json:"message"`
	Service                    string `json:"service"`
	ProjectID                  int    `json:"project_id"`
	Channel                    int    `json:"channel"`
	Status                     int    `json:"status"`
	Driver                     int    `json:"driver"`
	ContactInfo                string `json:"contact_info"`
}
