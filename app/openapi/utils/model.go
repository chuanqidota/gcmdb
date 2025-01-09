package utils

import (
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
)

type model struct {

}

var Model = new(model)


func (m *model)ModelAll()([]models.Model,error){
	_models := make([]models.Model, 0)
	if err := database.DB.Model(&models.Model{}).
		Scan(&_models).Error; err != nil {
		return nil,err
	}
	return _models,nil
}

