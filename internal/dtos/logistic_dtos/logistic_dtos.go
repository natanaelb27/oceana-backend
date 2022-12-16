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
	PreorderDay  uint       `json:"preorder_day"`
	IsPrintLabel bool       `json:"is_print_label"`
	Courier      []*Courier `json:"couriers"`
}

type CourierServiceCheck struct {
	ID uint `json:"id"`
}

type CourierCheck struct {
	ID             uint                   `json:"id"`
	CourierService []*CourierServiceCheck `json:"courier_services"`
}
type LogisticCheck struct {
	PreorderDay  uint            `json:"preorder_day"`
	IsPrintLabel bool            `json:"is_print_label"`
	Courier      []*CourierCheck `json:"couriers"`
}

func (source *Logistic) ToLogisticCheck() *LogisticCheck {
	courierCheck := make([]*CourierCheck, 0)
	for _, courier := range source.Courier {
		courierCheck = append(courierCheck, &CourierCheck{
			ID: courier.ID,
		})
	}
	return &LogisticCheck{
		PreorderDay:  source.PreorderDay,
		IsPrintLabel: source.IsPrintLabel,
		Courier:      courierCheck,
	}
}
