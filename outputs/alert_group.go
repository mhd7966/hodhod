package outputs

type AlertGroupOutput struct {
	Name      string   `json:"name"`
	UserID    int      `json:"user_id"`
	IsDefault bool     `json:"is_default"`
	Alerts    []string `json:"alerts"`
}
