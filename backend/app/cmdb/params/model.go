package params

type PatchModelGroupIdBody struct {
	GroupId uint `json:"group_id" binding:"required" label:"模型组ID"`
}

type PatchModelBody struct {
	Name        string `json:"name" binding:"required" label:"模型中文名"`
	Description string `json:"description" label:"描述"`
	IsUsable    bool   `json:"is_usable" label:"是否在用"`
	Icon        string `json:"icon" label:"图标"`
	Order       uint   `json:"order" label:"排序"`
}
