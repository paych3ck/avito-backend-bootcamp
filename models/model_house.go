package models

import (
	"time"
)

type House struct {
	Id        int32     `json:"id"`
	Address   string    `json:"address"`
	Year      int32     `json:"year"`
	Developer *string   `json:"developer,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdateAt  time.Time `json:"update_at,omitempty"`
}
