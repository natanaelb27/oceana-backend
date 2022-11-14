package dtos

type Category struct {
	ID       int         `json:"id"`
	Name     string      `json:"name"`
	Level    int         `json:"level"`
	Category []*Category `json:"category"`
}
