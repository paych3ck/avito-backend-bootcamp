package models

type DummyLoginGet500Response struct {
	Message   string `json:"message"`
	RequestId string `json:"request_id,omitempty"`
	Code      int32  `json:"code,omitempty"`
}
