package params

// UpdateModelFieldBody
// @Description: 修改模型字段 只能修改这些
type UpdateModelFieldBody struct {
	FieldGroupId uint   `json:"field_group_id"`
	Name         string `json:"name"`
	IsUnique     bool   `json:"is_unique"`
	IsRequired   bool   `json:"is_required"`
	Order        uint   `json:"order"`
}
