package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	in_dtos "git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

type RajaOngkirRepository interface {
	GetShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (result *in_dtos.GetRajaOngkirShippingCostResponse, err error)
}

type rajaOngkirRepository struct {
	client http.Client
}

type RajaOngkirRepositoryConfig struct {
	Client http.Client
}

func NewRajaOngkirRepository(c *RajaOngkirRepositoryConfig) RajaOngkirRepository {
	return &rajaOngkirRepository{
		client: c.Client,
	}
}

func (r *rajaOngkirRepository) GetShippingCost(request *dtos.GetRajaOngkirShippingCostRequest) (result *in_dtos.GetRajaOngkirShippingCostResponse, err error) {
	data, err := json.Marshal(request)

	if err != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	url := os.Getenv("RAJAONGKIR_API_URL")
	key := os.Getenv("RAJAONGKIR_API_KEY")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cost", url), bytes.NewBuffer(data))

	if err != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Key", key)
	res, err := r.client.Do(req)

	if err != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	if res.StatusCode != http.StatusOK {
		return nil, e.ErrDatabaseOperationFailed
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	err = json.Unmarshal(bodyBytes, &result)

	if err != nil {
		return nil, e.ErrDatabaseOperationFailed
	}

	return
}
