package models

type RegisterPostRequest struct {
	Email    string   `json:"email,omitempty"`
	Password string   `json:"password,omitempty"`
	UserType UserType `json:"user_type,omitempty"`
}
