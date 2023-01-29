package inputs

type NewContactProjectAlertGroup struct {
	ContactID int  `json:"contact_id" validate:"required"`
	Severity  int  `json:"severity" validate:"required,number,min=0,max=5"`
	Call      bool `json:"call" default:"false" validate:"required"`
	SMS       bool `json:"sms" default:"false" validate:"required"`
	Email     bool `json:"email" default:"false" validate:"required"`
	Webhook   bool `json:"webhook" default:"false" validate:"required"`
}

type UpdateContactProjectAlertGroup struct {
	Severity int  `json:"severity" validate:"number,min=0,max=5"`
	Call     bool `json:"call" default:"false"`
	SMS      bool `json:"sms" default:"false"`
	Email    bool `json:"email" default:"false"`
	Webhook  bool `json:"webhook" default:"false"`
}
