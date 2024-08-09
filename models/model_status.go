package models

type Status string

const (
	CREATED       Status = "created"
	APPROVED      Status = "approved"
	DECLINED      Status = "declined"
	ON_MODERATION Status = "on moderation"
)
