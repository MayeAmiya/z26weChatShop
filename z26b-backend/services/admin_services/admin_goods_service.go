package admin_services

import (
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AdminGoodsService struct {
	db *gorm.DB
}

func NewAdminGoodsService(db *gorm.DB) AdminGoodsServiceInterface {
	return &AdminGoodsService{db: db}
}

// GetProducts 获取商品列表（管理后台）
func (s *AdminGoodsService) GetProducts(page, pageSize int, keyword, categoryId, tagId, status string) ([]internal.SPU, int64, error) {
	var products []internal.SPU
	var total int64

	query := s.db.Model(&internal.SPU{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}
	if tagId != "" {
		query = query.Where("id IN (SELECT spu_id FROM spu_tag WHERE tag_id = ?)", tagId)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	// Load tags
	for i := range products {
		var spuTags []internal.SPUTag
		s.db.Preload("Tag").Where("spu_id = ?", products[i].ID).Find(&spuTags)
		var tags []internal.Tag
		for _, st := range spuTags {
			if st.Tag != nil {
				tags = append(tags, *st.Tag)
			}
		}
		products[i].Tags = tags
	}

	return products, total, nil
}

// GetProduct 获取商品详情（管理后台）
func (s *AdminGoodsService) GetProduct(id string) (*internal.SPU, []internal.SKU, error) {
	var product internal.SPU
	if err := s.db.Preload("Category").First(&product, "id = ?", id).Error; err != nil {
		return nil, nil, err
	}

	var spuTags []internal.SPUTag
	s.db.Preload("Tag").Where("spu_id = ?", id).Find(&spuTags)
	var tags []internal.Tag
	for _, st := range spuTags {
		if st.Tag != nil {
			tags = append(tags, *st.Tag)
		}
	}
	product.Tags = tags

	// 获取SKUs
	var skus []internal.SKU
	if err := s.db.Where("spu_id = ?", id).Find(&skus).Error; err != nil {
		return nil, nil, err
	}

	return &product, skus, nil
}

// UpdateProduct 更新商品
func (s *AdminGoodsService) UpdateProduct(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().Unix()
	return s.db.Model(&internal.SPU{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteProduct 删除商品
func (s *AdminGoodsService) DeleteProduct(id string) error {
	return s.db.Delete(&internal.SPU{}, "id = ?", id).Error
}

// GetCategories 获取分类列表（管理后台）
func (s *AdminGoodsService) GetCategories() ([]internal.Category, error) {
	var categories []internal.Category
	err := s.db.Order("sort ASC").Find(&categories).Error
	return categories, err
}

// CreateCategory 创建分类
func (s *AdminGoodsService) CreateCategory(category *internal.Category) error {
	return s.db.Create(category).Error
}

// UpdateCategory 更新分类
func (s *AdminGoodsService) UpdateCategory(id string, updates map[string]interface{}) error {
	return s.db.Model(&internal.Category{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory 删除分类
func (s *AdminGoodsService) DeleteCategory(id string) error {
	return s.db.Delete(&internal.Category{}, "id = ?", id).Error
}

// CreateProduct 创建商品
func (s *AdminGoodsService) CreateProduct(spu *internal.SPU, skus []internal.SKU) error {
	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建SPU
	now := time.Now().Unix()
	spu.CreatedAt = now
	spu.UpdatedAt = now
	if err := tx.Create(spu).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建SKU
	for i := range skus {
		skus[i].SPUID = spu.ID
		skus[i].CreatedAt = now
		skus[i].UpdatedAt = now
		if err := tx.Create(&skus[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// UpdateProductStatus 更新商品状态
func (s *AdminGoodsService) UpdateProductStatus(id, status string) error {
	return s.db.Model(&internal.SPU{}).Where("id = ?", id).Update("status", status).Error
}
