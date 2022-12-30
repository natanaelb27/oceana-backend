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
	LogoURL        string            `json:"logo_url"`
	IsSelected     bool              `json:"is_selected"`
	CourierService []*CourierService `json:"courier_services"`
}

type Logistic struct {
	IsEditPreorderDay bool       `json:"is_edit_preorder_day"`
	PreorderDay       uint       `json:"preorder_day"`
	IsPrintLabel      bool       `json:"is_print_label"`
	Courier           []*Courier `json:"couriers"`
}
