package services

import (
	"errors"
	"go_core/config"
	"go_core/models"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // 使用环境变量获取密钥

// Claims 是自定义的 JWT Claims 结构体
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// CreateUser 用于创建新用户
func CreateUser(user models.User) error {
	// 检查用户是否已存在
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	// 创建新用户
	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail 根据邮箱查找用户
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// CheckPassword 验证密码（实际项目中应该加密存储并验证）
func CheckPassword(storedPassword, providedPassword string) bool {
	// 简单示例，实际应该使用密码哈希进行比较（如 bcrypt）
	return storedPassword == providedPassword
}

// GenerateToken 生成 JWT Token
func GenerateToken(user models.User) (string, error) {
	// 设置 Token 过期时间为 24 小时
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 使用 Unix 时间戳表示过期时间
			Issuer:    "my-gin-project",      // 可以设置为应用名称
		},
	}

	// 创建 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名 Token
	return token.SignedString(jwtKey)
}

// ValidateToken 验证 JWT Token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证 Token 的签名方法是否为 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
