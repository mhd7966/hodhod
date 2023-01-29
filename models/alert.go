package models

import (
	"gorm.io/gorm"
)

type Alert struct {
	gorm.Model
	ID       int    `json:"id"`
	Message  string `json:"message"`
	Severity int    `json:"severity"`
	ParentID int    `json:"parent_id"`
}

type AlertGroup struct {
	gorm.Model
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	IsDefault bool   `json:"is_dafault"`
}

type AlertGroupAlerts struct {
	gorm.Model
	ID           int `json:"id"`
	AlertGroupID int `json:"alert_group_id"`
	AlertID      int `json:"alert_id"`
}

type ProjectAlertGroup struct {
	gorm.Model
	ID           int    `json:"id"`
	Name         string `json:"name"`
	UserID       int    `json:"user_id"`
	ProjectID    int    `json:"project_id"`
	Service      string `json:"service"`
	AlertGroupID int    `json:"alert_group_id"`
}

type ContactProjectAlertGroup struct {
	gorm.Model
	ID                  int  `json:"id"`
	ProjectAlertGroupID int  `json:"project_alert_group_id"`
	ContactID           int  `json:"contact_id"`
	Severity            int  `json:"severity"`
	Call                bool `json:"call"`
	SMS                 bool `json:"sms"`
	Email               bool `json:"email"`
	Webhook             bool `json:"webhook"`
}
