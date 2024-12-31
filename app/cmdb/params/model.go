package params

type PatchModelGroupIdBody struct {
	GroupId uint `json:"group_id"`
	ModelId uint `json:"model_id"`
}

type PatchModelBody struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsUsable    bool   `json:"is_usable"`
	Icon        string `json:"icon"`
	Order       uint   `json:"order"`
}
