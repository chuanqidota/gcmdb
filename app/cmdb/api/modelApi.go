package api

import "github.com/gin-gonic/gin"

type model struct {
}

var Model = new(model)

func (m *model) CreateModel(c *gin.Context) {
}

func (m *model) ListModel(c *gin.Context) {
}

func (m *model) RetrieveModel(c *gin.Context) {

}

func (m *model) UpdateModel(c *gin.Context) {

}

func (m *model) DeleteModel(c *gin.Context) {

}
