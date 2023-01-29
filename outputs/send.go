package outputs

type Send struct {
	Message string   `json:"message"`
	SMS     []string `json:"sms"`
	Call    []string `json:"call"`
}
