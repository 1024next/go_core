package services

import (
	"errors"
	"go_core/config"
	"go_core/models"
	"go_core/utils"

	"github.com/gin-gonic/gin"
)

// GetProductsWithPagination 获取产品列表，并返回分页信息
func GetProductsWithPagination(c *gin.Context) ([]models.Product, utils.Pagination, error) {
	// 获取分页参数
	pagination := utils.GetPagination(c)

	// 获取偏移量和限制
	offset, limit := pagination.Paginate()

	// 查询产品列表
	var products []models.Product
	if err := config.DB.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, pagination, err
	}

	// 查询总记录数
	var total int64
	if err := config.DB.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, pagination, err
	}

	// 设置总记录数
	pagination.Total = total

	return products, pagination, nil
}

func CreateProduct(product *models.Product) error {
	// 验证产品信息是否有效
	if product.Name == "" || product.Price <= 0 {
		return errors.New("Invalid product data")
	}

	// 创建产品
	if err := config.DB.Create(&product).Error; err != nil {
		return err
	}

	return nil
}
