package admin

import (
	admin_services "z26b-backend/services/admin_services"

	"gorm.io/gorm"
)

// Handler 管理后台处理器
type Handler struct {
	AdminGoodsService    admin_services.AdminGoodsServiceInterface
	AdminCategoryService admin_services.AdminCategoryServiceInterface
	DB                   *gorm.DB // 暂时保留，用于其他功能迁移
}

// NewHandler 创建处理器实例
func NewHandler(adminGoodsService admin_services.AdminGoodsServiceInterface, adminCategoryService admin_services.AdminCategoryServiceInterface, db *gorm.DB) *Handler {
	return &Handler{
		AdminGoodsService:    adminGoodsService,
		AdminCategoryService: adminCategoryService,
		DB:                   db,
	}
}
