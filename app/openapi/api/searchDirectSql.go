package api

import (
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type searchDirectSql struct {
}

var SearchDirectSql = new(searchDirectSql)

func (sds *searchDirectSql) Search(c *gin.Context) {
	uuid := c.Param("uuid")
	var searchDirectSql models.SearchDirectSql
	if err := database.DB.Model(&models.SearchDirectSql{}).
		Where(map[string]any{"uuid": uuid}).
		First(&searchDirectSql).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	sql := searchDirectSql.Sql
	params := searchDirectSql.Params
	var result interface{}
	if err := database.DB.Raw(sql, params).Scan(&result).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	response.Success(c, "查询成功", result)
}
