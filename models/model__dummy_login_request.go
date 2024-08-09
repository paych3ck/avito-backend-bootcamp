package models

type DummyLoginRequest struct {
	UserType UserType `json:"user_type" binding:"required"`
}
