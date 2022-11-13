package dtos

type Response struct {
	Code     int `json:"code"`
	Data     any `json:"data"`
	Error    any `json:"error"`
	Message  any `json:"message"`
	Metadata any `json:"metadata"`
}
