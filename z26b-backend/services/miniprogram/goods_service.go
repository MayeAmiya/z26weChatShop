package miniprogram

import (
	"z26b-backend/internal"

	"gorm.io/gorm"
)

// GoodsService 商品服务
type GoodsService struct {
	db *gorm.DB
}

func NewGoodsService(db *gorm.DB) GoodsServiceInterface {
	return &GoodsService{db: db}
}

// GetGoodsList 获取商品列表（小程序端）
func (s *GoodsService) GetGoodsList(page, pageSize int, categoryID, search string) ([]internal.SPU, int64, error) {
	var goods []internal.SPU
	query := s.db.Where("status = ?", "ENABLED")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if search != "" {
		query = query.Where("name LIKE ? OR detail LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Model(&internal.SPU{}).Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&goods).Error

	return goods, total, err
}

// GetGoodDetail 获取商品详情
func (s *GoodsService) GetGoodDetail(id string) (*internal.SPU, []internal.SKU, error) {
	var good internal.SPU
	if err := s.db.First(&good, "id = ?", id).Error; err != nil {
		return nil, nil, err
	}

	var skus []internal.SKU
	err := s.db.Where(`"SPUID" = ?`, id).Find(&skus).Error

	return &good, skus, err
}

// GetSKUDetail 获取SKU详情
func (s *GoodsService) GetSKUDetail(id string) (*internal.SKU, error) {
	var sku internal.SKU
	err := s.db.Preload("SPU").First(&sku, "id = ?", id).Error
	return &sku, err
}

// GetSKUsBySpuID 获取商品的所有SKU
func (s *GoodsService) GetSKUsBySpuID(spuID string) ([]internal.SKU, error) {
	var skus []internal.SKU
	err := s.db.Where(`"SPUID" = ?`, spuID).Find(&skus).Error
	return skus, err
}

// GetCategories 获取分类列表
func (s *GoodsService) GetCategories() ([]internal.Category, error) {
	var categories []internal.Category
	err := s.db.Order("sort ASC").Find(&categories).Error
	return categories, err
}

// SearchGoods 搜索商品
func (s *GoodsService) SearchGoods(keyword string, page, pageSize int) ([]internal.SPU, int64, error) {
	var goods []internal.SPU
	query := s.db.Where("status = ? AND (name LIKE ? OR detail LIKE ?)", "ENABLED", "%"+keyword+"%", "%"+keyword+"%")

	var total int64
	query.Model(&internal.SPU{}).Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Find(&goods).Error

	return goods, total, err
}

// GetHomeSwiper 获取首页轮播图
func (s *GoodsService) GetHomeSwiper() ([]internal.Swiper, error) {
	var swipers []internal.Swiper
	err := s.db.Order("priority DESC").Find(&swipers).Error
	return swipers, err
}

// GetHomeContent 获取首页富文本内容
func (s *GoodsService) GetHomeContent(key string) (*internal.HomeContent, error) {
	var content internal.HomeContent
	err := s.db.Where("key = ? AND enabled = ?", key, true).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

// GetHomeCategories 获取首页分类
func (s *GoodsService) GetHomeCategories() ([]internal.Category, error) {
	var categories []internal.Category
	err := s.db.Where("parent_id IS NULL OR parent_id = ''").Order("sort ASC").Find(&categories).Error
	return categories, err
}

// GetPromotions 获取促销活动
func (s *GoodsService) GetPromotions() ([]internal.Promotion, error) {
	var promotions []internal.Promotion
	err := s.db.Order("promotion_status DESC").Limit(10).Find(&promotions).Error
	return promotions, err
}
