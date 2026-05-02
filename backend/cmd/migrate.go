/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	authModels "gcmdb/app/auth/models"
	auditModels "gcmdb/app/audit/models"
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "迁移表",
	Long:  "自动迁移表 go run main.go migrate",
	Run: func(cmd *cobra.Command, args []string) {
		err := database.DB.AutoMigrate(
			&authModels.User{},
			&models.ModelGroup{},
			&models.Model{},
			&models.ModelRelation{},
			&models.ModelRelationType{},
			&models.ModelFieldGroup{},
			&models.ModelField{},
			&models.ModelFieldRelation{},
			&models.ModelFieldUnique{},
			&models.Instance{},
			&models.InstanceRelation{},
			&models.SearchDirectSql{},
			&auditModels.AuditLog{},
		)
		if err != nil {
			logger.Error("迁移出错:", err.Error())
			return
		}
		logger.Info("迁移成功")

		// 预置 admin 用户
		seedAdmin()
	},
}

func seedAdmin() {
	var count int64
	database.DB.Model(&authModels.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		logger.Info("admin 用户已存在，跳过创建")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("生成密码哈希失败:", err.Error())
		return
	}
	admin := authModels.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Token:        uuid.New().String(),
		IsActive:     true,
		IsAdmin:      true,
	}
	if err := database.DB.Create(&admin).Error; err != nil {
		logger.Error("创建 admin 用户失败:", err.Error())
		return
	}
	logger.Info("admin 用户创建成功, Token:", admin.Token)
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
