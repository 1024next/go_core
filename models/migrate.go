package models

import "go_core/config"

// Migrate 执行数据库迁移
func Migrate() {
	// 执行所有模型的迁移
	err := config.DB.AutoMigrate(
		&User{},
		&Product{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
