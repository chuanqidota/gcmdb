package api

import "github.com/gin-gonic/gin"

type modelGroup struct {
}

var ModelGroup = new(modelGroup)

// CreateModelGroup
//
//	@Description: 创建模型分组
//	@receiver m
//	@param c
func (m *modelGroup) CreateModelGroup(c *gin.Context) {

}

// ListModelGroup
//
//	@Description: 模型分组列表
//	@receiver m
//	@param c
func (m *modelGroup) ListModelGroup(c *gin.Context) {

}

// RetrieveModelGroup
//
//	@Description: 获取具体的一个模型分组
//	@receiver m
//	@param c
func (m *modelGroup) RetrieveModelGroup(c *gin.Context) {

}

// UpdateModelGroup
//
//	@Description: 更新模型分组
//	@receiver m
//	@param c
func (m *modelGroup) UpdateModelGroup(c *gin.Context) {
}

// DeleteModelGroup
//
//	@Description: 删除模型分组
//	@receiver m
//	@param c
func (m *modelGroup) DeleteModelGroup(c *gin.Context) {

}
