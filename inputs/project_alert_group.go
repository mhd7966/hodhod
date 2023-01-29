package inputs

type ProjectAlertGroupBody struct {
	Name         string `json:"name" validate:"required"`
	ProjectID    int    `json:"project_id" validate:"required"`
	Service      string `json:"service" validate:"required"`
}
