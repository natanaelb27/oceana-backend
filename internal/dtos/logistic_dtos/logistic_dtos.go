package logistics_dtos

type CourierService struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsSelected  bool   `json:"is_selected"`
}

type Courier struct {
	ID             uint              `json:"id"`
	Code           string            `json:"code"`
	Name           string            `json:"name"`
	IsSelected     bool              `json:"is_selected"`
	CourierService []*CourierService `json:"courier_services"`
}

type Logistics struct {
	PreorderDay  uint       `json:"preorder_day"`
	IsPrintLabel bool       `json:"is_print_label"`
	Courier      []*Courier `json:"courier"`
}
