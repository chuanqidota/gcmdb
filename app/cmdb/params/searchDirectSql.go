package params

import "gorm.io/datatypes"

type UpdateSearchDirectSql struct {
	Name   string         `json:"name"`
	Sql    string         `json:"sql"`
	Params datatypes.JSON `json:"params"`
}
