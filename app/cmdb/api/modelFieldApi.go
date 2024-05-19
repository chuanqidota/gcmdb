package api

import "github.com/gin-gonic/gin"

type modelField struct {
}

var ModelField = new(modelField)

func (m *modelField) CreateModelField(c *gin.Context) {
}

func (m *modelField) ListModelField(c *gin.Context) {
}

func (m *modelField) RetrieveModelField(c *gin.Context) {
}

func (m *modelField) UpdateModelField(c *gin.Context) {
}

func (m *modelField) DeleteModelField(c *gin.Context) {
}
