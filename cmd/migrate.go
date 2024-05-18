/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gcmdb/app/cmdb/models"
	"gcmdb/pkg/database"
	"gcmdb/pkg/logger"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "迁移表",
	Long:  "自动迁移表 go run main.go migrate",
	Run: func(cmd *cobra.Command, args []string) {
		err := database.DB.AutoMigrate(
			&models.ModelGroup{},
			&models.Model{},
			&models.ModelFieldGroup{},
			&models.ModelField{},
			&models.ModelRelation{},
			&models.ModelRelationType{},
			&models.Instance{},
			&models.InstanceRelation{},
		)
		if err != nil {
			logger.Error("迁移出错:", err.Error())
		} else {
			logger.Info("迁移成功")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
