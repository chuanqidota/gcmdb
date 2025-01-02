package params

import "gorm.io/datatypes"

type UpdateInstance struct {
	Data datatypes.JSON `json:"data"`
}
