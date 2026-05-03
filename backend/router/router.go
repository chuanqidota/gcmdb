package router

import (
	authApi "gcmdb/app/auth/api"
	auditApi "gcmdb/app/audit/api"
	cmdbApi "gcmdb/app/cmdb/api"
	openApi "gcmdb/app/openapi/api"
	taskApi "gcmdb/app/tasks/api"
	"gcmdb/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Engine() *gin.Engine {
	router := gin.Default()

	// 认证路由
	auth := router.Group("v1/auth")
	{
		auth.POST("login", authApi.Login.Login)
		auth.POST("logout", authApi.Login.Logout)
		auth.GET("me", middleware.SessionAuthMiddleware(), authApi.Login.Me)
		auth.POST("change-password", middleware.SessionAuthMiddleware(), authApi.Login.ChangePassword)
		auth.POST("reset-password", middleware.SessionAuthMiddleware(), authApi.Login.ResetPassword)

		// 用户管理（管理员）
		auth.GET("users", middleware.SessionAuthMiddleware(), authApi.User.ListUsers)
		auth.POST("users", middleware.SessionAuthMiddleware(), authApi.User.CreateUser)
		auth.PATCH("users/:id", middleware.SessionAuthMiddleware(), authApi.User.PatchUser)
	}

	v1 := router.Group("v1")

	// CMDB 页面接口 — CORS + session 认证 + 审计
	cmdb := v1.Group("cmdb").Use(middleware.CORSMiddleware()).Use(middleware.SessionAuthMiddleware()).Use(middleware.AuditMiddleware())
	{
		// 模型分组
		cmdb.POST("models-group", cmdbApi.ModelGroup.CreateModelGroup)       // 创建模型分组
		cmdb.GET("models-group", cmdbApi.ModelGroup.ListModelGroup)          // 模型分组查询
		cmdb.GET("models-group/:id", cmdbApi.ModelGroup.RetrieveModelGroup)  //  模型分组详情
		cmdb.PATCH("models-group/:id", cmdbApi.ModelGroup.PatchModelGroup)   // 修改模型分组
		cmdb.DELETE("models-group/:id", cmdbApi.ModelGroup.DeleteModelGroup) // 删除模型分组

		// 模型
		cmdb.POST("models", cmdbApi.Model.CreateModel)                     // 创建模型
		cmdb.GET("models", cmdbApi.Model.ListModel)                        // 模型查询
		cmdb.GET("models/:id", cmdbApi.Model.RetrieveModel)                // 模型详情
		cmdb.PUT("models/:id", cmdbApi.Model.UpdateModel)                  // 修改模型
		cmdb.PATCH("models/:id", cmdbApi.Model.PatchModelGroupId)          // 修改模型分组
		cmdb.DELETE("models/:id", cmdbApi.Model.DeleteModel)               // 删除模型

		// 模型关系
		cmdb.POST("models-relation", cmdbApi.ModelRelation.CreateModelRelation)       // 创建模型关系
		cmdb.GET("models-relation", cmdbApi.ModelRelation.ListModelRelation)          // 模型关系查询
		cmdb.DELETE("models-relation/:id", cmdbApi.ModelRelation.DeleteModelRelation) // 删除模型关系

		// 模型关系类型
		cmdb.POST("models-relation-type", cmdbApi.ModelRelationType.CreateModelRelationType)    // 创建模型关系类型
		cmdb.GET("models-relation-type", cmdbApi.ModelRelationType.ListModelRelationType)       // 模型关系类型查询
		cmdb.PUT("models-relation-type/:id", cmdbApi.ModelRelationType.UpdateModelRelationType) // 修改模型关系类型
		cmdb.DELETE("models-relation-type/:id", cmdbApi.ModelRelationType.DeleteModelRelationType)  // 删除模型关系类型

		// 模型字段分组
		cmdb.POST("models-field-group", cmdbApi.ModelFieldGroup.CreateModelFieldGroup)       // 创建模型字段分组
		cmdb.GET("models-field-group/:id", cmdbApi.ModelFieldGroup.RetrieveModelFieldGroup)  // 查询模型字段分组详情
		cmdb.PUT("models-field-group/:id", cmdbApi.ModelFieldGroup.UpdateModelFieldGroup)    // 修改模型字段分组
		cmdb.DELETE("models-field-group/:id", cmdbApi.ModelFieldGroup.DeleteModelFieldGroup) // 删除模型字段分组

		// 模型唯一字段
		cmdb.POST("models-field-unique", cmdbApi.ModelFieldUnique.CreateModelFieldUnique)            // 创建模型唯一字段
		cmdb.GET("models-field-unique/:model_id", cmdbApi.ModelFieldUnique.RetrieveModelFieldUnique) // 展示模型的唯一字段
		cmdb.DELETE("models-field-unique/:id", cmdbApi.ModelFieldUnique.DeleteModelFieldUnique)      // 删除模型唯一字段

		// 模型字段
		cmdb.POST("models-field", cmdbApi.ModelField.CreateModelField)            // 创建模型字段
		cmdb.GET("models-field/:model_id", cmdbApi.ModelField.RetrieveModelField) // 查询模型字段详情
		cmdb.PUT("models-field/:id", cmdbApi.ModelField.UpdateModelField)         // 修改模型字段
		cmdb.DELETE("models-field/:id", cmdbApi.ModelField.DeleteModelField)      // 删除模型字段

		// 模型字段关系
		cmdb.POST("models-field-relation", cmdbApi.ModelFieldRelation.CreateModelFieldRelation)               // 创建模型字段关系
		cmdb.GET("models-field-relation/:source_model_id", cmdbApi.ModelFieldRelation.ListModelFieldRelation) // 展示模型字段关系
		cmdb.DELETE("models-field-relation/:id", cmdbApi.ModelFieldRelation.DeleteModelFieldRelation)         // 删除模型字段关系

		// 实例
		cmdb.POST("instance", cmdbApi.Instance.CreateInstance)           // 创建实例
		cmdb.GET("instances/:model_id", cmdbApi.Instance.ListInstance)   // 查询实例列表
		cmdb.GET("instance/:id", cmdbApi.Instance.RetrieveInstance)      // 查询实例详情
		cmdb.PUT("instance/:id", cmdbApi.Instance.UpdateInstance)        // 更新实例
		cmdb.DELETE("instance/:id", cmdbApi.Instance.DeleteInstance)     // 删除实例

		// 实例关系
		cmdb.POST("instance-relation", cmdbApi.InstanceRelation.CreateInstanceRelation)                      // 创建实例关系
		cmdb.GET("source-target-relation/:source_id", cmdbApi.InstanceRelation.SourceTargetInstanceRelation) // 源到目标的实例关系
		cmdb.GET("target-source-relation/:target_id", cmdbApi.InstanceRelation.TargetSourceInstanceRelation) // 目标到源的实例关系
		cmdb.DELETE("instance-relation/:id", cmdbApi.InstanceRelation.DeleteInstanceRelation)                // 删除实例关系
		cmdb.DELETE("instance-relation-by-keys", cmdbApi.InstanceRelation.DeleteInstanceRelationByKeys)    // 通过源/目标删除实例关系

		// 直接查询sql
		cmdb.POST("search-direct-sql", cmdbApi.SearchDirectSql.CreateSearchDirectSql)       // 创建直接查询sql
		cmdb.GET("search-direct-sql", cmdbApi.SearchDirectSql.ListSearchDirectSql)          // 查询直接查询sql
		cmdb.GET("search-direct-sql/:id", cmdbApi.SearchDirectSql.ExecuteSearchDirectSql)   // 执行直接查询sql
		cmdb.PUT("search-direct-sql/:id", cmdbApi.SearchDirectSql.UpdateSearchDirectSql)    // 修改直接查询sql
		cmdb.DELETE("search-direct-sql/:id", cmdbApi.SearchDirectSql.DeleteSearchDirectSql) // 删除直接查询sql

		// 审计日志
		cmdb.GET("audit-log", auditApi.Audit.ListAuditLog)                                    // 查询审计日志
		cmdb.GET("audit-log/:resource_type/:resource_id", auditApi.Audit.RetrieveAuditLog)    // 查询资源审计历史

		// 任务
		cmdb.POST("sync-instance-relation", taskApi.InstanceRelation.SyncInstanceRelation) // 同步实例关系
	}

	// 对外开放接口 — session 或 token 认证
	openapi := router.Group("openapi").Use(middleware.OpenAPIAuthMiddleware())
	{
		openapi.GET("model/:range", openApi.Model.ModelRange)                                      // 模型操作
		openapi.GET("instance/:action", openApi.Instance.InstanceQuery)                            // 实例查询
		openapi.POST("instance/:action", openApi.Instance.InstanceWrite)                           // 实例写入
		openapi.POST("instance-relation/:action", openApi.InstanceRelation.InstanceRelationAction) // 实例关系操作
	}

	return router
}
