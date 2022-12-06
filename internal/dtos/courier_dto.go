package dtos

type CourierService struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Courier struct {
	ID             uint              `json:"id"`
	Code           string            `json:"code"`
	Name           string            `json:"name"`
	CourierService []*CourierService `json:"courier_services"`
}
