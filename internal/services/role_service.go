package services

import (
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	r "git.garena.com/azmi-adhani/group-project-backend/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleService interface {
	WithTx(tx *gorm.DB) *roleService
	CronRolesCount() error
	Create(role *models.Role) (*models.Role, error)
	GetAll() ([]*models.Role, error)
	GetById(id uuid.UUID) (*models.Role, error)
	Update(id uuid.UUID, newModel *models.Role) (*models.Role, error)
	Delete(id uuid.UUID) (*models.Role, error)
	MultipleOperationFailingWithTransaction(model *models.Role) (*models.Role, error)
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

// WithTrx enables repository with transaction
func (u *roleService) WithTx(tx *gorm.DB) *roleService {
	u.roleRepository = u.roleRepository.WithTx(tx)
	return u
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

func (p *roleService) MultipleOperationFailingWithTransaction(model *models.Role) (*models.Role, error) {
	_, err := p.roleRepository.Create(model)
	if err != nil {
		return nil, err
	}

	// failing create
	_, err = p.roleRepository.CreateFailing(nil)
	if err != nil {
		return nil, err
	}

	result, err := p.roleRepository.Create(model)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (p *roleService) CronRolesCount() error {
	roles, err := p.roleRepository.GetAll()
	if err != nil {
		return err
	}
	fmt.Println("CRON: roles count is", len(roles))
	return nil
}
