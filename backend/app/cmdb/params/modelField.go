package params

// UpdateModelFieldBody
// @Description: 修改模型字段 只能修改这些
type UpdateModelFieldBody struct {
	FieldGroupId *uint  `json:"field_group_id" label:"字段分组ID"`
	Name         string `json:"name" label:"字段名称" binding:"required"`
	IsRequired   bool   `json:"is_required" label:"是否必须"`
	IsIndexed    bool   `json:"is_indexed" label:"是否索引"`
	Order        uint   `json:"order" label:"排序"`
}
