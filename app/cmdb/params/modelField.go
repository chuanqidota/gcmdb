package params

type UpdateModelFieldBody struct {
	FieldGroupId uint   `json:"field_group_id"`
	Name         string `json:"name"`
	IsUnique     bool   `json:"is_unique"`
	IsRequired   bool   `json:"is_required"`
	Order        uint   `json:"order"`
}
