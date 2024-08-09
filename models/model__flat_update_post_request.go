package models

type FlatUpdatePostRequest struct {
	Id     int32  `json:"id"`
	Status Status `json:"status,omitempty"`
}
