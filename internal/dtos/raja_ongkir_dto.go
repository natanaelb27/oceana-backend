package dtos

type RajaOngkirShippingServiceCost struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type RajaOngkirShippingService struct {
	Service     string                          `json:"service"`
	Description string                          `json:"description"`
	Cost        []RajaOngkirShippingServiceCost `json:"cost"`
}

type RajaOngkirShippingCostStatus struct {
	Code uint `json:"code"`
}

type RajaOngkirShippingCostResult struct {
	Code  string                      `json:"code"`
	Name  string                      `json:"name"`
	Costs []RajaOngkirShippingService `json:"costs"`
}

type RajaOngkirShippingCostResponse struct {
	Status  RajaOngkirShippingCostStatus   `json:"status"`
	Results []RajaOngkirShippingCostResult `json:"results"`
}

type GetRajaOngkirShippingCostResponse struct {
	RajaOngkir RajaOngkirShippingCostResponse `json:"rajaongkir"`
}
