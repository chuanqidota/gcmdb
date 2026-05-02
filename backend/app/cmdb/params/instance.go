package params

import "gorm.io/datatypes"

type UpdateInstance struct {
	Data datatypes.JSON `json:"data" binding:"required"`
}

type ListInstance struct {
	Search  string `form:"search" label:"模糊搜索"`
	Field   string `form:"field" label:"指定字段"`
	Value   string `form:"value" label:"指定值"`
	Compare string `form:"compare" label:"比较符"`
	Limit   int    `form:"limit" label:"分页limit"`
	Offset  int    `form:"offset" label:"分页offset"`
}

func ValidatorCompare(compare string) bool {
	switch compare {
	case "=", "like", "!=":
		return true
	}
	return false
}
