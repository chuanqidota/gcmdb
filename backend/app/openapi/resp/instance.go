package resp

import "gorm.io/datatypes"

type FulltextInstance struct {
	ID        uint           `json:"id" label:"实例ID"`
	ModelId   uint           `json:"model_id" label:"模型ID"`
	ModelName string         `json:"model_name" label:"模型名称"`
	Data      datatypes.JSON `json:"data" label:"实例数据"`
}
