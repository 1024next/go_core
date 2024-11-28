package handlers

import (
	"fmt"
	"game_service/config"
	"game_service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 创建游戏
func CreateGame(c *gin.Context) {
	var game models.Game
	// 绑定 JSON 数据到 game 结构体
	if err := c.ShouldBindJSON(&game); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// 校验游戏数据
	if game.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game name is required"})
		return
	}

	// 保存游戏到数据库
	if err := config.DB.Create(&game).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game", "details": err.Error()})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusCreated, gin.H{
		"message": "Game created successfully",
		"game":    game,
	})
}

// 获取所有游戏
func GetGames(c *gin.Context) {
	var games []models.Game

	// 查询游戏，并计算每个游戏的房间数量
	result := config.DB.Table("games").Select("games.id, games.name, games.status, COUNT(rooms.id) AS room_count").
		// 使用 LEFT JOIN 以确保即使没有房间也能返回游戏
		Joins("LEFT JOIN rooms ON rooms.game_id = games.id").
		Group("games.id").Find(&games)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 返回游戏和房间数量
	c.JSON(http.StatusOK, games)
}

// 根据ID获取游戏
func GetGameByID(c *gin.Context) {
	id := c.Param("game_id")
	fmt.Print(c)
	// 将ID转换为整型
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的游戏ID"})
		return
	}

	var game models.Game
	// 查询指定ID的游戏
	if err := config.DB.First(&game, gameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "游戏不存在"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// 删除游戏
func DeleteGame(c *gin.Context) {
	id := c.Param("game_id")

	// 将ID转换为整型
	gameID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的游戏ID"})
		return
	}

	// 删除游戏
	if err := config.DB.Delete(&models.Game{}, gameID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除游戏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "游戏已删除"})
}
