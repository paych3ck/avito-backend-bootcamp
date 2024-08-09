package models

type UserType string

const (
	CLIENT    UserType = "client"
	MODERATOR UserType = "moderator"
)
