package api

import (
	"errors"
	"fmt"
	"gcmdb/app/cmdb/models"
	"gcmdb/app/cmdb/params"
	"gcmdb/pkg/database"
	"gcmdb/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type searchDirectSql struct {
}

var SearchDirectSql = new(searchDirectSql)

// CreateSearchDirectSql
//
//	@Description: 创建直接查询sql
//	@receiver sds
//	@param c
func (sds *searchDirectSql) CreateSearchDirectSql(c *gin.Context) {
	var body models.SearchDirectSql
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	body.Uuid = strings.Replace(uuid.New().String(), "-", "", -1)
	if err := database.DB.Model(&models.SearchDirectSql{}).Create(&body).Error; err != nil {
		response.Fail(c, fmt.Sprintf("创建失败-%s", err.Error()))
		return
	}
	response.Success(c, "创建成功", nil)
}

// ListSearchDirectSql
//
//	@Description: 查询直接查询sql
//	@receiver sds
//	@param c
func (sds *searchDirectSql) ListSearchDirectSql(c *gin.Context) {
	var query params.CommonQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	search := query.Search
	limit, offset := query.Limit, query.Offset
	if limit == 0 {
		limit = 10
	}
	db := database.DB.Model(&models.SearchDirectSql{})
	if search != "" {
		db.Where("name like ?", "%"+search+"%").
			Or("uuid like ?", "%"+search+"%").
			Or("sql like ?", "%"+search+"%")
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	var searchDirectSqls []models.SearchDirectSql
	if err := db.Limit(limit).Offset(offset).Scan(&searchDirectSqls).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	results := params.CommonList{
		Count:   count,
		Results: searchDirectSqls,
	}
	response.Success(c, "执行成功", results)
}

// ExecuteSearchDirectSql
//
//	@Description: 执行查询
//	@receiver sds
//	@param c
func (sds *searchDirectSql) ExecuteSearchDirectSql(c *gin.Context) {
	id := c.Param("id")
	var searchDirectSql models.SearchDirectSql
	if err := database.DB.Model(&models.SearchDirectSql{}).
		Where(map[string]any{"id": id}).
		First(&searchDirectSql).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
	sql := searchDirectSql.Sql
	if sql == "" {
		response.Fail(c, "sql语句为空")
		return
	}
	var result interface{}
	if err := database.DB.Raw(sql).Scan(&result).Error; err != nil {
		response.Fail(c, fmt.Sprintf("查询失败-%s", err.Error()))
		return
	}
}

// UpdateSearchDirectSql
//
//	@Description: 更新直接查询sql
//	@receiver sds
//	@param c
func (sds *searchDirectSql) UpdateSearchDirectSql(c *gin.Context) {
	id := c.Param("id")
	var body params.UpdateSearchDirectSql
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, fmt.Sprintf("参数错误-%s", err.Error()))
		return
	}
	data := map[string]any{
		"created_at": time.Now(),
		"name":       body.Name,
		"sql":        body.Sql,
	}
	if err := database.DB.Model(&models.SearchDirectSql{}).
		Where(map[string]any{"id": id}).
		Updates(data).Error; err != nil {
		response.Fail(c, fmt.Sprintf("更新失败-%s", err.Error()))
		return
	}
	response.Success(c, "执行成功", nil)

}

// DeleteSearchDirectSql
//
//	@Description: 删除直接查询sql
//	@receiver sds
//	@param c
func (sds *searchDirectSql) DeleteSearchDirectSql(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Model(&models.SearchDirectSql{}).
		Where(map[string]any{"id": id}).
		Delete(&models.SearchDirectSql{}).Error; err != nil {
		response.Fail(c, fmt.Sprintf("删除失败-%s", err.Error()))
		return
	}
	response.Success(c, "删除成功", nil)
}
