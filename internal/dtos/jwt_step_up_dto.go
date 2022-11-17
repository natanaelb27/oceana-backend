package dtos

type JwtStepUp struct {
	Scope string `json:"scope"`
}

func NewJwtStepUp(scope string) *JwtStepUp {
	return &JwtStepUp{
		Scope: scope,
	}
}
