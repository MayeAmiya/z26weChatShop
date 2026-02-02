package admin_services

import (
	"errors"
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AdminOrderService struct {
	db *gorm.DB
}

func NewAdminOrderService(db *gorm.DB) AdminOrderServiceInterface {
	return &AdminOrderService{db: db}
}

// GetOrders 获取订单列表（管理后台）
func (s *AdminOrderService) GetOrders(page, pageSize int, keyword, status, userID string, startDate, endDate *time.Time) ([]internal.Order, int64, error) {
	var orders []internal.Order
	var total int64

	query := s.db.Model(&internal.Order{}).Preload("OrderItems").Preload("OrderItems.SKU").Preload("OrderItems.SKU.SPU")

	if keyword != "" {
		query = query.Where("order_no LIKE ? OR receiver_name LIKE ? OR receiver_phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetOrderByID 根据ID获取订单详情
func (s *AdminOrderService) GetOrderByID(orderID string) (*internal.Order, error) {
	var order internal.Order
	if err := s.db.Preload("OrderItems").Preload("OrderItems.SKU").Preload("OrderItems.SKU.SPU").
		Preload("User").Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订单不存在")
		}
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (s *AdminOrderService) UpdateOrderStatus(orderID, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// 如果是发货状态，设置发货时间
	if status == "shipped" {
		updates["shipped_at"] = time.Now()
	}

	// 如果是完成状态，设置完成时间
	if status == "completed" {
		updates["completed_at"] = time.Now()
	}

	if err := s.db.Model(&internal.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// UpdateOrderShipping 更新订单物流信息
func (s *AdminOrderService) UpdateOrderShipping(orderID, shippingCompany, trackingNumber string) error {
	updates := map[string]interface{}{
		"shipping_company": shippingCompany,
		"tracking_number":  trackingNumber,
		"updated_at":       time.Now(),
	}

	if err := s.db.Model(&internal.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// GetOrderStats 获取订单统计信息
func (s *AdminOrderService) GetOrderStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总订单数
	var totalOrders int64
	if err := s.db.Model(&internal.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// 待付款订单数
	var pendingPayment int64
	if err := s.db.Model(&internal.Order{}).Where("status = ?", "pending_payment").Count(&pendingPayment).Error; err != nil {
		return nil, err
	}
	stats["pending_payment"] = pendingPayment

	// 待发货订单数
	var pendingShipment int64
	if err := s.db.Model(&internal.Order{}).Where("status = ?", "pending_shipment").Count(&pendingShipment).Error; err != nil {
		return nil, err
	}
	stats["pending_shipment"] = pendingShipment

	// 已发货订单数
	var shipped int64
	if err := s.db.Model(&internal.Order{}).Where("status = ?", "shipped").Count(&shipped).Error; err != nil {
		return nil, err
	}
	stats["shipped"] = shipped

	// 已完成订单数
	var completed int64
	if err := s.db.Model(&internal.Order{}).Where("status = ?", "completed").Count(&completed).Error; err != nil {
		return nil, err
	}
	stats["completed"] = completed

	// 已取消订单数
	var cancelled int64
	if err := s.db.Model(&internal.Order{}).Where("status = ?", "cancelled").Count(&cancelled).Error; err != nil {
		return nil, err
	}
	stats["cancelled"] = cancelled

	// 总销售额
	var totalSales float64
	if err := s.db.Model(&internal.Order{}).Where("status IN ?", []string{"TO_RECEIVE", "FINISHED"}).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalSales).Error; err != nil {
		return nil, err
	}
	stats["total_sales"] = totalSales

	// 本月销售额
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var monthlySales float64
	if err := s.db.Model(&internal.Order{}).Where("status IN ? AND created_at >= ?", []string{"TO_RECEIVE", "FINISHED"}, startOfMonth).
		Select("COALESCE(SUM(final_price), 0)").Scan(&monthlySales).Error; err != nil {
		return nil, err
	}
	stats["monthly_sales"] = monthlySales

	return stats, nil
}

// RefundOrder 退款订单
func (s *AdminOrderService) RefundOrder(orderID string, refundAmount float64, refundReason string) error {
	// 获取订单信息
	var order internal.Order
	if err := s.db.Where("id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("订单不存在")
		}
		return err
	}

	// 检查订单状态
	if order.Status != "shipped" && order.Status != "completed" {
		return errors.New("只有已发货或已完成的订单才能退款")
	}

	// 检查退款金额
	if refundAmount > order.FinalPrice {
		return errors.New("退款金额不能超过订单总金额")
	}

	// 更新订单状态为退款中
	updates := map[string]interface{}{
		"status":        "refunding",
		"refund_amount": refundAmount,
		"refund_reason": refundReason,
		"updated_at":    time.Now(),
	}

	if err := s.db.Model(&internal.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

// ConfirmRefund 确认退款
func (s *AdminOrderService) ConfirmRefund(orderID string) error {
	updates := map[string]interface{}{
		"status":     "refunded",
		"updated_at": time.Now(),
	}

	if err := s.db.Model(&internal.Order{}).Where("id = ?", orderID).Updates(updates).Error; err != nil {
		return err
	}

	return nil
}