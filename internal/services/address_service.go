package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
)

type AddressService interface {
	GetProfileAddressess(userID uuid.UUID) (model []*models.Address, err error)
	Create(address *models.Address) (*models.Address, error)
	GetProfileAddressById(id uuid.UUID) (model *models.Address, err error)
	UpdateProfileAddress(id uuid.UUID, newModel *models.Address) (model *models.Address, err error)
}

type addressService struct {
	addressRepository r.AddressRepository
}

type AddressServiceConfig struct {
	AddressRepository r.AddressRepository
}

func NewAddressService(c *AddressServiceConfig) AddressService {
	return &addressService{
		addressRepository: c.AddressRepository,
	}
}

func (a *addressService) Create(model *models.Address) (*models.Address, error) {
	isAvailable, _ := a.addressRepository.GetProfileAddressess(model.UserID)
	if len(isAvailable) == 0 {
		model.IsDefault = true
	} else {
		model.IsDefault = false
	}

	return a.addressRepository.Create(model)
}

func (a *addressService) GetProfileAddressess(userID uuid.UUID) ([]*models.Address, error) {
	return a.addressRepository.GetProfileAddressess(userID)
}

func (a *addressService) GetProfileAddressById(id uuid.UUID) (model *models.Address, err error) {
	return a.addressRepository.GetProfileAddressById(id)
}

func (a *addressService) UpdateProfileAddress(id uuid.UUID, newModel *models.Address) (model *models.Address, err error) {
	return a.addressRepository.UpdateProfileAddress(id, newModel)
}
