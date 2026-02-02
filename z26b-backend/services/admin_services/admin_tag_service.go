package admin_services

import (
	"errors"
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AdminTagService struct {
	db *gorm.DB
}

func NewAdminTagService(db *gorm.DB) AdminTagServiceInterface {
	return &AdminTagService{db: db}
}

// GetTags 获取标签列表（管理后台）
func (s *AdminTagService) GetTags(page, pageSize int, keyword, status string) ([]internal.Tag, int64, error) {
	var tags []internal.Tag
	var total int64

	query := s.db.Model(&internal.Tag{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("sort_order ASC, created_at DESC").Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetTagByID 根据ID获取标签详情
func (s *AdminTagService) GetTagByID(tagID string) (*internal.Tag, error) {
	var tag internal.Tag
	if err := s.db.Where("id = ?", tagID).First(&tag).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("标签不存在")
		}
		return nil, err
	}
	return &tag, nil
}

// CreateTag 创建标签
func (s *AdminTagService) CreateTag(name, description, color string, sortOrder int) (*internal.Tag, error) {
	// 检查标签名是否已存在
	var existingTag internal.Tag
	if err := s.db.Where("name = ?", name).First(&existingTag).Error; err == nil {
		return nil, errors.New("标签名已存在")
	}

	tag := internal.Tag{
		Name:        name,
		Description: description,
		Color:       color,
		SortOrder:   sortOrder,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// UpdateTag 更新标签
func (s *AdminTagService) UpdateTag(tagID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	if err := s.db.Model(&internal.Tag{}).Where("id = ?", tagID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTag 删除标签
func (s *AdminTagService) DeleteTag(tagID string) error {
	// 检查是否有商品使用此标签
	var count int64
	if err := s.db.Model(&internal.SPU{}).Where("FIND_IN_SET(?, tags)", tagID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("该标签正在被商品使用，无法删除")
	}

	if err := s.db.Where("id = ?", tagID).Delete(&internal.Tag{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTagStatus 更新标签状态
func (s *AdminTagService) UpdateTagStatus(tagID, status string) error {
	if err := s.db.Model(&internal.Tag{}).Where("id = ?", tagID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

// GetAllActiveTags 获取所有激活的标签
func (s *AdminTagService) GetAllActiveTags() ([]internal.Tag, error) {
	var tags []internal.Tag
	if err := s.db.Where("status = ?", "active").Order("sort_order ASC, created_at DESC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}