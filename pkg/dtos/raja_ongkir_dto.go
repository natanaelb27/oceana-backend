package dtos

import (
	in_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
)

type GetRajaOngkirShippingCostRequest struct {
	Origin      uint   `json:"origin" binding:"required"`
	Destination uint   `json:"destination" binding:"required"`
	Weight      uint   `json:"weight" binding:"required"`
	Courier     string `json:"courier" binding:"required"`
	ProductIDs  []uint `json:"product_ids" binding:"required"`
}

type RajaOngkirShippingServiceCost struct {
	Value int    `json:"value"`
	Etd   string `json:"etd"`
	Note  string `json:"note"`
}

type RajaOngkirShippingService struct {
	CourierServiceID uint                            `json:"courier_service_id"`
	Service          string                          `json:"service"`
	Description      string                          `json:"description"`
	Cost             []RajaOngkirShippingServiceCost `json:"cost"`
}

type GetRajaOngkirShippingCostResponse struct {
	Code  string                      `json:"code"`
	Name  string                      `json:"name"`
	Costs []RajaOngkirShippingService `json:"costs"`
}

func NewGetRajaOngkirShippingCostResponse(source *in_dtos.GetRajaOngkirShippingCostResponse) (target *GetRajaOngkirShippingCostResponse) {
	services := []RajaOngkirShippingService{}

	for _, service := range source.RajaOngkir.Results[0].Costs {
		costs := []RajaOngkirShippingServiceCost{}

		for _, cost := range service.Cost {
			costs = append(costs, RajaOngkirShippingServiceCost{
				Value: cost.Value,
				Etd:   cost.Etd,
				Note:  cost.Note,
			})
		}

		services = append(services, RajaOngkirShippingService{
			Service:     service.Service,
			Description: service.Description,
			Cost:        costs,
		})
	}

	target = &GetRajaOngkirShippingCostResponse{
		Code:  source.RajaOngkir.Results[0].Code,
		Name:  source.RajaOngkir.Results[0].Name,
		Costs: services,
	}

	return
}
