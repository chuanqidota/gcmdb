package params

import "gorm.io/datatypes"

type UpdateInstance struct {
	Data datatypes.JSON `json:"data"`
}

type ListInstance struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Field    string `form:"field"`
	Value    string `form:"value"`
	Category string `form:"category"`
}
