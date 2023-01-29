package models

type AuthResponse struct{
	UserID int `json:"id" mapstructure:"id"`
	Name string `json:"name" mapstructure:"name"`
	Email string` json:"email" mapstructure:"email"`
}