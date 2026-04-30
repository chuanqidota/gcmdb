package params

type PatchModelGroupBody struct {
	Name        string `json:"name" label:"模型分组名称" binding:"required"`
	Description string `json:"description" label:"描述" binding:"required"`
}
