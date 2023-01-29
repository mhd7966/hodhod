package outputs

type LogInfo struct {
	ID                int    `json:"id"`
	Contact           string `json:"contact"`
	AlertGroup        string `json:"alert_group"`
	ProjectAlertGroup string `json:"project_alert_group"`
	Message           string `json:"message"`
	Service           string `json:"service"`
	Channel           int    `json:"channel"`
	Status            int    `json:"status"`
	ContactInfo       string `json:"contact_info"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
