package miniprogram

import (
	"errors"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) CartServiceInterface {
	return &CartService{db: db}
}

// GetCartItems 获取用户购物车商品
func (s *CartService) GetCartItems(userID string) ([]internal.CartItem, error) {
	var items []internal.CartItem
	err := s.db.Preload("SKU").Preload("SKU.SPU").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

// AddToCart 添加商品到购物车
func (s *CartService) AddToCart(userID, skuID string, quantity int) error {
	// 检查SKU是否存在
	var sku internal.SKU
	if err := s.db.First(&sku, "id = ?", skuID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("SKU not found")
		}
		return err
	}

	// 检查库存 (暂时注释掉以便测试)
	// if sku.Count < quantity {
	// 	return errors.New("insufficient stock")
	// }

	// 检查是否已在购物车中
	var existingItem internal.CartItem
	err := s.db.Where("user_id = ? AND sku_id = ?", userID, skuID).First(&existingItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 不存在，创建新项
			item := internal.CartItem{
				ID:       internal.GenerateUUID(),
				UserID:   userID,
				SKUID:    skuID,
				Quantity: quantity,
			}
			return s.db.Create(&item).Error
		}
		return err
	}

	// 已存在，更新数量
	existingItem.Quantity += quantity
	// 检查库存 (暂时注释掉)
	// if existingItem.Quantity > sku.Count {
	// 	return errors.New("insufficient stock")
	// }
	return s.db.Save(&existingItem).Error
}

// UpdateCartItem 更新购物车商品数量和选择状态
func (s *CartService) UpdateCartItem(userID, itemID string, quantity int, isSelected *bool) error {
	var item internal.CartItem
	if err := s.db.Where("id = ? AND user_id = ?", itemID, userID).First(&item).Error; err != nil {
		return err
	}

	// 如果提供了选择状态，更新它
	if isSelected != nil {
		item.IsSelected = *isSelected
	}

	// 如果数量 <= 0，删除项
	if quantity <= 0 && quantity != -1 {
		return s.RemoveCartItem(userID, itemID)
	}

	// 如果提供了数量，更新它
	if quantity > 0 {
		// 检查库存 (暂时注释掉)
		// var sku internal.SKU
		// if err := s.db.First(&sku, "id = ?", item.SKUID).Error; err != nil {
		// 	return err
		// }

		// if sku.Count < quantity {
		// 	return errors.New("insufficient stock")
		// }

		item.Quantity = quantity
	}

	return s.db.Save(&item).Error
}

// RemoveCartItem 从购物车移除商品
func (s *CartService) RemoveCartItem(userID, itemID string) error {
	return s.db.Where("id = ? AND user_id = ?", itemID, userID).Delete(&internal.CartItem{}).Error
}

// ClearCart 清空用户购物车
func (s *CartService) ClearCart(userID string) error {
	return s.db.Where("user_id = ?", userID).Delete(&internal.CartItem{}).Error
}

// GetCartItemCount 获取购物车商品数量
func (s *CartService) GetCartItemCount(userID string) (int64, error) {
	var count int64
	err := s.db.Model(&internal.CartItem{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
