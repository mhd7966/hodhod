package inputs

type AlertGroupBody struct {
	Name      string `json:"name" validate:"required"`
	AlertIDs  []int  `json:"alert_ids" validate:"required"`
	IsDefault bool   `json:"is_dafault" default:"false"`
}

type AlertGroupUpdate struct {
	Name string `json:"name" validate:"required"`
}

type AlertGroupAlerts struct {
	AlertIDs []int `json:"alert_ids" validate:"required"`
}
