package models

type LoginPostRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
