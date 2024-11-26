package routes

import (
	"go_core/controllers"
	"go_core/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Logger())
	// Public routes
	r.POST("/api/register", controllers.RegisterUser)
	r.POST("/api/login", controllers.LoginUser)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.GET("/products", controllers.GetProducts)
		protected.POST("/products", controllers.CreateProduct)
		protected.POST("/upload", controllers.UploadFile)
	}

	return r
}
