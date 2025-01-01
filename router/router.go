package router

import (
	"gcmdb/app/cmdb/api"
	"gcmdb/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	router := gin.Default()
	gcmdb := router.Group("gcmdb")
	v1 := gcmdb.Group("v1").Use(middleware.CORSMiddleware())

	{
		// 模型分组
		v1.POST("model-group", api.ModelGroup.CreateModelGroup)       // 创建模型分组
		v1.GET("model-group", api.ModelGroup.ListModelGroup)          // 模型分组查询
		v1.GET("model-group/:id", api.ModelGroup.RetrieveModelGroup)  //  模型分组详情
		v1.PATCH("model-group/:id", api.ModelGroup.PatchModelGroup)   // 修改模型分组
		v1.DELETE("model-group/:id", api.ModelGroup.DeleteModelGroup) // 删除模型分组

		// 模型
		v1.POST("model", api.Model.CreateModel)                     // 创建模型
		v1.GET("model", api.Model.ListModel)                        // 模型查询
		v1.GET("model/:id", api.Model.RetrieveModel)                // 模型详情
		v1.PUT("model/:id", api.Model.UpdateModel)                  // 修改模型
		v1.PATCH("model/change-group", api.Model.PatchModelGroupId) // 修改模型分组
		v1.DELETE("model/:id", api.Model.DeleteModel)               // 删除模型

		// 模型关系
		v1.POST("model-relation", api.ModelRelation.CreateModelRelation)       // 创建模型关系
		v1.GET("model-relation", api.ModelRelation.ListModelRelation)          // 模型关系查询
		v1.DELETE("model-relation/:id", api.ModelRelation.DeleteModelRelation) // 删除模型关系

		// 模型关系类型
		v1.POST("model-relation-type", api.ModelRelationType.CreateModelRelationType)    // 创建模型关系类型
		v1.GET("model-relation-type", api.ModelRelationType.ListModelRelationType)       // 模型关系类型查询
		v1.PUT("model-relation-type/:id", api.ModelRelationType.UpdateModelRelationType) // 修改模型关系类型
		v1.DELETE("model-relation-type", api.ModelRelationType.DeleteModelRelationType)  // 删除模型关系类型

		// 模型字段分组
		v1.POST("model-field-group", api.ModelFieldGroup.CreateModelFieldGroup)       // 创建模型字段分组
		v1.GET("model-field-group/:id", api.ModelFieldGroup.RetrieveModelFieldGroup)  // 查询模型字段分组详情
		v1.PUT("model-field-group/:id", api.ModelFieldGroup.UpdateModelFieldGroup)    // 修改模型字段分组
		v1.DELETE("model-field-group/:id", api.ModelFieldGroup.DeleteModelFieldGroup) // 删除模型字段分组

		// 模型唯一字段
		v1.POST("model-field-unique", api.ModelFieldUnique.CreateModelFieldUnique)            // 模型唯一字段
		v1.GET("model-field-unique/:model_id", api.ModelFieldUnique.RetrieveModelFieldUnique) // 展示模型的唯一字段
		v1.DELETE("model-field-unique/:id", api.ModelFieldUnique.DeleteModelFieldUnique)      // 删除模型唯一字段

		// 模型字段
		v1.POST("model-field", api.ModelField.CreateModelField)            // 创建模型字段-todo 实例里面异步补充字段
		v1.GET("model-field/:model_id", api.ModelField.RetrieveModelField) // 查询模型字段详情
		v1.PUT("model-field/:id", api.ModelField.UpdateModelField)         // 修改模型字段
		v1.DELETE("model-field/:id", api.ModelField.DeleteModelField)      // 删除模型字段-todo 实例里面删除字段

		// 模型字段关系
		v1.POST("model-field-relation", api.ModelFieldRelation.CreateModelFieldRelation)               // 创建模型字段关系-todo 异步创建实例
		v1.GET("model-field-relation/:source_model_id", api.ModelFieldRelation.ListModelFieldRelation) // 展示模型字段关系
		v1.DELETE("model-field-relation/:id", api.ModelFieldRelation.DeleteModelFieldRelation)         // 删除模型字段关系-todo 需要操作实例

		// 实例
		v1.POST("instance", api.Instance.CreateInstance)       // 创建实例- todo 未做唯一性判断
		v1.GET("instance", api.Instance.ListInstance)          // 查询实例-todo
		v1.GET("instance/:id", api.Instance.RetrieveInstance)  // 查询实例详情
		v1.PUT("instance/:id", api.Instance.UpdateInstance)    // 更新实例-todo
		v1.DELETE("instance/:id", api.Instance.DeleteInstance) // 删除实例

		// 实例关系
		v1.POST("instance-relation", api.InstanceRelation.CreateInstanceRelation)       // 创建实例关系
		v1.GET("instance-relation", api.InstanceRelation.ListInstanceRelation)          // 查询实例关系 - todo
		v1.DELETE("instance-relation/:id", api.InstanceRelation.DeleteInstanceRelation) // 删除实例关系

	}

	return router
}
