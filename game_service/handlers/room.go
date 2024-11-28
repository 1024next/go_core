package handlers

import (
	"encoding/json"
	"fmt"
	"game_service/config"
	"game_service/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取房间中的所有玩家
func GetRoomPlayers(c *gin.Context) {
	gameID := c.Param("game_id")
	roomID := c.Param("room_id")

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	var room models.Room
	// 查询房间
	if err := config.DB.First(&room, "id = ? AND game_id = ?", roomID, validatedGameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在或不属于该游戏"})
		return
	}

	// 解析 Players 字段，获取玩家列表
	var players []string
	if room.Players == "" {
		players = []string{} // 初始化为空切片
	}

	// 返回玩家列表
	c.JSON(http.StatusOK, gin.H{"room_id": room.ID, "players": players})
}

// 验证游戏是否存在并返回其 uint 类型的 ID
func validateGame(gameID string) (uint, error) {
	id, err := strconv.Atoi(gameID) // 将字符串转换为整数
	if err != nil {
		return 0, err
	}

	var game models.Game
	if err := config.DB.First(&game, id).Error; err != nil {
		return 0, err
	}

	return uint(id), nil
}

// 获取指定游戏的房间列表
func GetRooms(c *gin.Context) {
	gameID := c.Param("game_id")

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	var rooms []models.Room
	// 查询指定游戏的房间
	if err := config.DB.Where("game_id = ?", validatedGameID).Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取房间列表失败"})
		return
	}

	c.JSON(http.StatusOK, rooms)
}

// 创建房间
func CreateRoom(c *gin.Context) {
	gameID := c.Param("game_id")

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置房间的 GameID
	room.GameID = validatedGameID

	// 保存房间
	if err := config.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建房间失败"})
		return
	}

	// 更新游戏的房间数量
	if err := updateGameRoomCount(validatedGameID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新游戏房间数量失败"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// 更新游戏的房间数量
func updateGameRoomCount(gameID uint) error {
	var game models.Game
	// 获取游戏并更新房间数量
	if err := config.DB.First(&game, gameID).Error; err != nil {
		return fmt.Errorf("游戏不存在: %v", err)
	}

	// 更新房间数量
	game.RoomCount++ // 假设数据库中有一个 `RoomCount` 字段
	if err := config.DB.Save(&game).Error; err != nil {
		return fmt.Errorf("更新游戏房间数量失败: %v", err)
	}

	return nil
}

// 用户加入房间
func JoinRoom(c *gin.Context) {
	gameID := c.Param("game_id")
	roomID := c.Param("room_id")

	var user struct {
		Username string `json:"username"` // 用户名
	}

	// 获取用户请求中的用户名
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户信息"})
		return
	}

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	var room models.Room
	// 查询房间
	if err := config.DB.First(&room, "id = ? AND game_id = ?", roomID, validatedGameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在或不属于该游戏"})
		return
	}

	// 解析 Players 字段，获取当前玩家列表
	var players []string

	if room.Players == "" {
		players = []string{} // 初始化为空切片
	} else {
		if err := json.Unmarshal([]byte(room.Players), &players); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析玩家列表失败"})
			return
		}

		// 检查玩家是否已经在房间中
		for _, player := range players {
			if player == user.Username {
				c.JSON(http.StatusConflict, gin.H{"error": "玩家已在房间中"})
				return
			}
		}

		// 检查房间是否已满
		if len(players) >= room.MaxSeats {
			c.JSON(http.StatusBadRequest, gin.H{"error": "房间已满"})
			return
		}

	}

	// 将用户添加到玩家列表
	players = append(players, user.Username)

	// 更新房间的 Players 字段
	updatedPlayers, err := json.Marshal(players)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新玩家列表失败"})
		return
	}

	room.Players = string(updatedPlayers)

	// 保存更新后的房间数据
	if err := config.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存房间信息失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "玩家成功加入房间",
		"room_id": room.ID,
		"players": players,
	})
}

// 用户退出房间
func LeaveRoom(c *gin.Context) {
	gameID := c.Param("game_id")
	roomID := c.Param("room_id")

	// 获取退出房间的用户信息
	var user struct {
		Username string `json:"username"` // 用户名
	}

	// 如果请求体中没有传递用户名，则返回错误
	if err := c.ShouldBindJSON(&user); err != nil || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体中缺少用户名信息"})
		return
	}

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	var room models.Room
	// 查询房间
	if err := config.DB.First(&room, "id = ? AND game_id = ?", roomID, validatedGameID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "房间不存在或不属于该游戏"})
		return
	}

	// 解析 Players 字段，获取当前玩家列表
	var players []string
	if room.Players == "" {
		players = []string{} // 如果 Players 字段为空，初始化为空切片
	} else {
		if err := json.Unmarshal([]byte(room.Players), &players); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "解析玩家列表失败"})
			return
		}
	}

	// 检查玩家是否在房间内
	playerFound := false
	for i, player := range players {
		if player == user.Username {
			// 找到玩家，移除
			players = append(players[:i], players[i+1:]...) // 从列表中移除玩家
			playerFound = true
			break
		}
	}

	if !playerFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "玩家不在房间中"})
		return
	}

	// 更新房间的 Players 字段
	updatedPlayers, err := json.Marshal(players)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新玩家列表失败"})
		return
	}

	room.Players = string(updatedPlayers)

	// 保存更新后的房间数据
	if err := config.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存房间信息失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "玩家成功退出房间",
		"room_id": room.ID,
		"players": players,
	})
}

// 删除房间
func DeleteRoom(c *gin.Context) {
	gameID := c.Param("game_id")
	roomID := c.Param("room_id")

	// 验证游戏是否存在
	validatedGameID, err := validateGame(gameID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效或不存在的游戏ID"})
		return
	}

	// 删除房间
	if err := config.DB.Where("id = ? AND game_id = ?", roomID, validatedGameID).Delete(&models.Room{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除房间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "房间已删除"})
}
