package miniprogram

import (
	"encoding/json"
	"errors"
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderServiceInterface {
	return &OrderService{db: db}
}

// GetOrderList 获取用户订单列表
func (s *OrderService) GetOrderList(userID, status string, page, pageSize int) ([]internal.Order, int64, error) {
	var orders []internal.Order
	var total int64

	query := s.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Model(&internal.Order{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Items.SKU").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error
	return orders, total, err
}

// GetOrderDetail 获取订单详情
func (s *OrderService) GetOrderDetail(orderID, userID string) (*internal.Order, error) {
	var order internal.Order
	err := s.db.Preload("Items.SKU.SPU").Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	return &order, err
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID string, items []internal.OrderItem, addressID string) (*internal.Order, error) {
	// 验证地址 (暂时注释掉)
	// var address internal.Address
	// if err := s.db.First(&address, "id = ? AND user_id = ?", addressID, userID).Error; err != nil {
	// 	return nil, errors.New("invalid address")
	// }

	// 创建一个默认地址用于测试 (暂时不需要)
	// address := internal.Address{...}

	// 计算总价
	var totalPrice float64
	for i := range items {
		var sku internal.SKU
		if err := s.db.First(&sku, "id = ?", items[i].SKUID).Error; err != nil {
			return nil, errors.New("invalid sku: " + items[i].SKUID)
		}
		// 检查库存 (暂时注释掉)
		// if sku.Count < items[i].Quantity {
		// 	return nil, errors.New("insufficient stock for sku: " + items[i].SKUID)
		// }
		items[i].Price = sku.Price
		if items[i].Price == 0 {
			items[i].Price = 1.0 // 默认价格为1，用于测试
		}
		totalPrice += items[i].Price * float64(items[i].Quantity)
	}

	// 创建订单
	deliveryInfo := map[string]interface{}{
		"name":          "Test User",
		"phone":         "1234567890",
		"countryName":   "China",
		"provinceName":  "Test Province",
		"cityName":      "Test City",
		"districtName":  "Test District",
		"detailAddress": "Test Address",
	}
	deliveryJSON, _ := json.Marshal(deliveryInfo)
	order := internal.Order{
		ID:           internal.GenerateUUID(),
		UserID:       userID,
		Status:       internal.OrderStatusToSend, // 直接设置为待发货，跳过支付
		DeliveryInfo: deliveryJSON,
		TotalPrice:   totalPrice,
		FinalPrice:   totalPrice,
		// Items:        items, // 移除，因为单独创建
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 设置 items 的 OrderID
	for i := range items {
		items[i].OrderID = order.ID
	}

	// 创建 order items
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 清空购物车中对应的商品
	for _, item := range items {
		tx.Where("user_id = ? AND sku_id = ?", userID, item.SKUID).Delete(&internal.CartItem{})
	}

	tx.Commit()

	// 重新获取订单，确保包含完整的items和关联数据
	var createdOrder internal.Order
	s.db.Preload("Items.SKU.SPU").First(&createdOrder, "id = ?", order.ID)

	return &createdOrder, nil
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID, userID, status string) error {
	return s.db.Model(&internal.Order{}).Where("id = ? AND user_id = ?", orderID, userID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// CancelOrder 取消订单 - 允许待支付和待发货状态的订单取消
func (s *OrderService) CancelOrder(orderID, userID string) error {
	// 允许取消的状态：待支付、待发货
	allowedStatuses := []string{internal.OrderStatusToPay, internal.OrderStatusToSend}

	result := s.db.Model(&internal.Order{}).
		Where("id = ? AND user_id = ? AND status IN ?", orderID, userID, allowedStatuses).
		Updates(map[string]interface{}{
			"status":     internal.OrderStatusCanceled,
			"updated_at": time.Now(),
		})

	if result.RowsAffected == 0 {
		return errors.New("订单不存在或无法取消")
	}
	return result.Error
}

// GetAdminOrderList 管理员获取订单列表
func (s *OrderService) GetAdminOrderList(status, orderNo, userID string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var orders []internal.Order
	var total int64

	query := s.db.Model(&internal.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderNo != "" {
		query = query.Where("id LIKE ?", "%"+orderNo+"%")
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Preload("Items.SKU.SPU").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	// 添加用户信息
	var result []map[string]interface{}
	for _, order := range orders {
		orderMap := map[string]interface{}{
			"order": order,
		}

		var user internal.User
		if err := s.db.First(&user, "id = ?", order.UserID).Error; err == nil {
			orderMap["user"] = user
		}

		result = append(result, orderMap)
	}

	return result, total, nil
}

// UpdateAdminOrderStatus 管理员更新订单状态
func (s *OrderService) UpdateAdminOrderStatus(orderID, status string) error {
	return s.db.Model(&internal.Order{}).Where("id = ?", orderID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
