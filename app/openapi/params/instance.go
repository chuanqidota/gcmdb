package params

import (
	"gorm.io/datatypes"
)

type CreateInstance struct {
	ModelAlias string `json:"model_alias"`
	Data datatypes.JSON `json:"data"`
}
