package models

type HouseCreatePostRequest struct {
	Address   string  `json:"address"`
	Year      int32   `json:"year"`
	Developer *string `json:"developer,omitempty"`
}
