package params

import (
	"gorm.io/datatypes"
)

type CreateInstance struct {
	ModelAlias string         `json:"model_alias" label:"模型别名" binding:"required"`
	Data       datatypes.JSON `json:"data" label:"实例数据" binding:"required"`
}

type UpdateInstance struct {
	Id   uint           `json:"id" label:"实例ID" binding:"required"`
	Data datatypes.JSON `json:"data" label:"更新数据" binding:"required"`
}

type DeleteInstance struct {
	ModelAlias string `json:"model_alias" label:"模型英文名" binding:"required"`
	Id         uint   `json:"id" label:"实例ID" binding:"required"`
}

type MulDeleteInstance struct {
	ModelAlias string `json:"model_alias" label:"模型英文名" binding:"required"`
	Ids        []uint `json:"ids" label:"实例ID" binding:"required"`
}

type DirectSearch struct {
	Uuid string `json:"uuid" label:"uuid" binding:"required"`
}

type FulltextInstance struct {
	Search     string `json:"search" binding:"required" label:"模糊搜索" `
	Offset     int    `form:"offset" label:"分页offset"`
	Limit      int    `form:"limit" label:"分页limit"`
	ModelAlias string `json:"model_alias" label:"模型英文名,中间用逗号分割"`
}

type Condition struct {
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
	Order  []string         `json:"order"`
	Where  []map[string]any `json:"where"`
}

type SearchInstance struct {
	Model     string    `json:"model" binding:"required" label:"模型英文名"`
	Fields    []string  `json:"fields" label:"查询字段"`
	Condition Condition `json:"__condition" label:"查询条件"`
}

type SourceTargetInstance struct {
	Id    int64  `form:"id" binding:"required" label:"实例id"`
	Model string `form:"model" binding:"required" label:"源/目标模型英文名称"`
}
