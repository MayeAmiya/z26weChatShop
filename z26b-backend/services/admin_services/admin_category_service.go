package admin_services

import (
	"fmt"
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AdminCategoryService struct {
	db *gorm.DB
}

func NewAdminCategoryService(db *gorm.DB) AdminCategoryServiceInterface {
	return &AdminCategoryService{db: db}
}

// GetCategories 获取分类列表
func (s *AdminCategoryService) GetCategories() ([]internal.Category, error) {
	var categories []internal.Category
	err := s.db.Order("sort ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

// GetCategory 获取分类详情
func (s *AdminCategoryService) GetCategory(id string) (*internal.Category, error) {
	var category internal.Category
	err := s.db.First(&category, "id = ?", id).Error
	return &category, err
}

// CreateCategory 创建分类
func (s *AdminCategoryService) CreateCategory(req *internal.Category) error {
	req.ID = internal.GenerateUUID()
	return s.db.Create(req).Error
}

// UpdateCategory 更新分类
func (s *AdminCategoryService) UpdateCategory(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.db.Model(&internal.Category{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory 删除分类
func (s *AdminCategoryService) DeleteCategory(id string) error {
	// 检查是否有商品使用此分类
	var count int64
	s.db.Model(&internal.SPU{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("分类下有商品，无法删除")
	}

	return s.db.Delete(&internal.Category{}, "id = ?", id).Error
}
