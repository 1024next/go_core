package main

import (
	"game_service/config"
	"game_service/models"
	"game_service/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，使用默认环境变量")
	}

	// 初始化数据库
	config.ConnectDB()

	// 自动迁移数据库
	if err := config.DB.AutoMigrate(&models.Room{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	// 在 ConnectDB() 函数中添加：
	if err := config.DB.AutoMigrate(&models.Game{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 设置路由并启动服务
	router := routes.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
