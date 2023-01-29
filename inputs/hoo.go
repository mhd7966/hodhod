package inputs

type Hoo struct {
	UserID    int      `json:"user_id" validate:"required"`
	AlertID   int      `json:"alert_id" validate:"required"`
	Values    []string `json:"values" validate:"required"`
	Service   string   `json:"service" validate:"required"`
	ProjectID int      `json:"project_id" validate:"required"`
}
