package services

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
)

type RoleService interface {
	Create(role *models.Role) (*models.Role, error)
	GetAll() ([]*models.Role, error)
	GetById(id uuid.UUID) (*models.Role, error)
	Update(id uuid.UUID, newModel *models.Role) (*models.Role, error)
	Delete(id uuid.UUID) (*models.Role, error)
}

type roleService struct {
	roleRepository r.RoleRepository
}

type RoleServiceConfig struct {
	RoleRepository r.RoleRepository
}

func NewRoleService(c *RoleServiceConfig) RoleService {
	return &roleService{
		roleRepository: c.RoleRepository,
	}
}

func (p *roleService) Create(model *models.Role) (*models.Role, error) {
	return p.roleRepository.Create(model)
}

func (p *roleService) GetAll() ([]*models.Role, error) {
	return p.roleRepository.GetAll()
}

func (p *roleService) GetById(id uuid.UUID) (*models.Role, error) {
	return p.roleRepository.GetById(id)
}

func (p *roleService) Update(id uuid.UUID, newModel *models.Role) (*models.Role, error) {
	return p.roleRepository.Update(id, newModel)
}

func (p *roleService) Delete(id uuid.UUID) (*models.Role, error) {
	return p.roleRepository.Delete(id)
}
