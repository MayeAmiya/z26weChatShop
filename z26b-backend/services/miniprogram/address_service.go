package miniprogram

import (
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AddressService struct {
	db *gorm.DB
}

func NewAddressService(db *gorm.DB) AddressServiceInterface {
	return &AddressService{db: db}
}

// GetAddressList 获取用户地址列表
func (s *AddressService) GetAddressList(userID string) ([]internal.Address, error) {
	var addresses []internal.Address
	err := s.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

// GetAddress 获取单个地址
func (s *AddressService) GetAddress(id, userID string) (*internal.Address, error) {
	var address internal.Address
	err := s.db.First(&address, "id = ? AND user_id = ?", id, userID).Error
	return &address, err
}

// CreateAddress 创建地址
func (s *AddressService) CreateAddress(address *internal.Address) error {
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now
	return s.db.Create(address).Error
}

// UpdateAddress 更新地址
func (s *AddressService) UpdateAddress(id, userID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.db.Model(&internal.Address{}).Where("id = ? AND user_id = ?", id, userID).Updates(updates).Error
}

// DeleteAddress 删除地址
func (s *AddressService) DeleteAddress(id, userID string) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&internal.Address{}).Error
}

// SetDefaultAddress 设置默认地址
func (s *AddressService) SetDefaultAddress(id, userID string) error {
	// 先取消所有默认地址
	if err := s.db.Model(&internal.Address{}).Where("user_id = ?", userID).Update("is_default", 0).Error; err != nil {
		return err
	}

	// 设置新的默认地址
	return s.db.Model(&internal.Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", 1).Error
}
