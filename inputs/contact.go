package inputs

type ContactBody struct {
	Name        string    `json:"name" validate:"required"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
}
