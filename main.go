package main

import (
	"go_core/config"
	"go_core/models"
	"go_core/routes"
)

func main() {
	// 初始化数据库
	config.InitDB()

	// 自动迁移
	models.Migrate()

	// 初始化路由
	r := routes.SetupRouter()

	// 运行服务
	r.Run(":8080")
}
