package controllers

import (
	"go_core/models"
	"go_core/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProducts 获取产品列表（带分页）
func GetProducts(c *gin.Context) {
	// 调用服务层获取分页产品列表
	products, pagination, err := services.GetProductsWithPagination(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching products"})
		return
	}

	response := models.NewSuccessResponse(gin.H{
		"list":       products,
		"pagination": pagination,
	})
	// 返回产品数据和分页信息
	c.JSON(http.StatusOK, response)
}

// CreateProduct 创建新产品
func CreateProduct(c *gin.Context) {
	var product models.Product
	// 将请求体中的数据绑定到 product 结构体
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// 调用服务层创建产品
	if err := services.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// 返回创建成功的响应
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully", "product": product})
}
