package dtos

type Category struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	Level    int         `json:"level"`
	Category []*Category `json:"category"`
}
