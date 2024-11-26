package controllers

import (
	"go_core/models"
	"go_core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册接口
func RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// 调用服务层创建用户
	if err := services.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login 用户登录接口，生成 JWT Token
func LoginUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// 查找用户
	dbUser, err := services.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// 验证密码
	if !services.CheckPassword(dbUser.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// 生成 Token
	token, err := services.GenerateToken(*dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate token"})
		return
	}

	// 返回 token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
