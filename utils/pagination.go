package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination 用于封装分页的参数
type Pagination struct {
	Page     int
	PageSize int
	Total    int64
}

// DefaultPagination 设置分页的默认值
func DefaultPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 10,
	}
}

// GetPagination 从请求中获取分页参数
func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 返回分页结构体
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Paginate 计算分页的偏移量
func (p *Pagination) Paginate() (int, int) {
	offset := (p.Page - 1) * p.PageSize
	limit := p.PageSize
	return offset, limit
}
