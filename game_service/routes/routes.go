package routes

import (
	"game_service/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 游戏路由
	gameRoutes := router.Group("/games")
	{
		gameRoutes.GET("/", handlers.GetGames)              // 获取所有游戏
		gameRoutes.POST("/", handlers.CreateGame)           // 创建游戏
		gameRoutes.GET("/:game_id", handlers.GetGameByID)   // 获取指定游戏
		gameRoutes.DELETE("/:game_id", handlers.DeleteGame) // 删除游戏

		// 游戏下的房间子路由，确保 game_id 路径参数不与游戏的 id 路径参数冲突
		roomRoutes := gameRoutes.Group("/:game_id/rooms")
		{
			roomRoutes.GET("/", handlers.GetRooms)                       // 获取指定游戏的房间列表
			roomRoutes.POST("/", handlers.CreateRoom)                    // 创建房间
			roomRoutes.POST("/:room_id/join", handlers.JoinRoom)         // 加入房间
			roomRoutes.DELETE("/:room_id/players", handlers.LeaveRoom)   // 退出房间
			roomRoutes.DELETE("/:room_id", handlers.DeleteRoom)          // 删除房间
			roomRoutes.GET("/:room_id/players", handlers.GetRoomPlayers) // 获取房间的玩家列表
		}
	}

	// WebSocket 路由
	router.GET("/ws", handlers.WebSocketHandler) // WebSocket 路由
	return router
}
