package models

type FlatCreatePostRequest struct {
	HouseId    int32 `json:"house_id"`
	FlatNumber int32 `json:"flat_number"`
	Price      int32 `json:"price"`
	Rooms      int32 `json:"rooms,omitempty"`
}
