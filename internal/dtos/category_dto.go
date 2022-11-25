package dtos

type Category struct {
	ID       uint        `json:"id"`
	Image    string      `json:"image"`
	Name     string      `json:"name"`
	Level    int         `json:"level"`
	Category []*Category `json:"category"`
}
