package inputs

type AlertBody struct {
	Message  string `json:"message" validate:"required"`
	Severity int    `json:"severity" validate:"number,min=0,max=5"`
	ParentID int    `json:"parent_id"`
}

type AlertUpdate struct {
	Message  string `json:"message"`
	Severity int    `json:"severity" validate:"number,min=0,max=5"`
}
