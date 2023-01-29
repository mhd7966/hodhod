package models

import "gorm.io/gorm"

type SendList struct {
	SMSList  []string `json:"sms_ist"`
	CallList []string `json:"call_list"`
}

type HooQuery struct {
	gorm.Model
	ID                  int  `json:"id"`
	ProjectAlertGroupID int  `json:"project_alert_group_id"`
	ContactID           int  `json:"contact_id"`
	Severity            int  `json:"severity"`
	Call                bool `json:"call"`
	SMS                 bool `json:"sms"`
	Email               bool `json:"email"`
	Webhook             bool `json:"webhook"`
	AlertGroupID        int  `json:"alert_group_id"`
}
