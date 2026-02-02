package admin_services

import (
	"time"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type AdminStatsService struct {
	db *gorm.DB
}

func NewAdminStatsService(db *gorm.DB) AdminStatsServiceInterface {
	return &AdminStatsService{db: db}
}

// GetStats 获取仪表板统计数据
func (s *AdminStatsService) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 用户统计
	var totalUsers int64
	if err := s.db.Model(&internal.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	stats["total_users"] = totalUsers

	// 新增用户（今日）- 使用时间戳
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayTimestamp := todayStart.Unix() * 1000 // 毫秒时间戳

	var newUsersToday int64
	if err := s.db.Model(&internal.User{}).Where("created_at >= ?", todayStart).Count(&newUsersToday).Error; err != nil {
		return nil, err
	}
	stats["new_users_today"] = newUsersToday

	// 新增用户（本周）
	weekAgo := now.AddDate(0, 0, -7)
	var newUsersWeek int64
	if err := s.db.Model(&internal.User{}).Where("created_at >= ?", weekAgo).Count(&newUsersWeek).Error; err != nil {
		return nil, err
	}
	stats["new_users_week"] = newUsersWeek

	// 商品统计
	var totalProducts int64
	if err := s.db.Model(&internal.SPU{}).Count(&totalProducts).Error; err != nil {
		return nil, err
	}
	stats["total_products"] = totalProducts

	// 在售商品
	var activeProducts int64
	if err := s.db.Model(&internal.SPU{}).Where("status = ?", "active").Count(&activeProducts).Error; err != nil {
		return nil, err
	}
	stats["active_products"] = activeProducts

	// 订单统计
	var totalOrders int64
	if err := s.db.Model(&internal.Order{}).Count(&totalOrders).Error; err != nil {
		return nil, err
	}
	stats["total_orders"] = totalOrders

	// 今日订单 - 使用毫秒时间戳（Order.CreatedAt 是 int64 毫秒时间戳）
	var ordersToday int64
	if err := s.db.Model(&internal.Order{}).Where("created_at >= ?", todayTimestamp).Count(&ordersToday).Error; err != nil {
		return nil, err
	}
	stats["orders_today"] = ordersToday

	// 本周订单
	weekAgoTimestamp := weekAgo.Unix() * 1000
	var ordersWeek int64
	if err := s.db.Model(&internal.Order{}).Where("created_at >= ?", weekAgoTimestamp).Count(&ordersWeek).Error; err != nil {
		return nil, err
	}
	stats["orders_week"] = ordersWeek

	// 有效订单状态（已发货、已完成等）
	validStatuses := []string{
		internal.OrderStatusToSend,    // 待发货
		internal.OrderStatusToReceive, // 待收货
		internal.OrderStatusFinished,  // 已完成
	}

	// 销售额统计 - 使用 final_price 字段
	var totalSales float64
	if err := s.db.Model(&internal.Order{}).Where("status IN ?", validStatuses).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalSales).Error; err != nil {
		return nil, err
	}
	stats["total_sales"] = totalSales

	// 今日销售额
	var salesToday float64
	if err := s.db.Model(&internal.Order{}).Where("status IN ? AND created_at >= ?", validStatuses, todayTimestamp).
		Select("COALESCE(SUM(final_price), 0)").Scan(&salesToday).Error; err != nil {
		return nil, err
	}
	stats["sales_today"] = salesToday

	// 本周销售额
	var salesWeek float64
	if err := s.db.Model(&internal.Order{}).Where("status IN ? AND created_at >= ?", validStatuses, weekAgoTimestamp).
		Select("COALESCE(SUM(final_price), 0)").Scan(&salesWeek).Error; err != nil {
		return nil, err
	}
	stats["sales_week"] = salesWeek

	// 分类统计
	var totalCategories int64
	if err := s.db.Model(&internal.Category{}).Count(&totalCategories).Error; err != nil {
		return nil, err
	}
	stats["total_categories"] = totalCategories

	// 评论统计
	var totalComments int64
	if err := s.db.Model(&internal.Comment{}).Count(&totalComments).Error; err != nil {
		return nil, err
	}
	stats["total_comments"] = totalComments

	return stats, nil
}

// GetSalesStats 获取销售额图表数据
func (s *AdminStatsService) GetSalesStats(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 有效订单状态
	validStatuses := []string{
		internal.OrderStatusToSend,
		internal.OrderStatusToReceive,
		internal.OrderStatusFinished,
	}

	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

		// 转换为毫秒时间戳
		startTimestamp := startOfDay.Unix() * 1000
		endTimestamp := endOfDay.Unix() * 1000

		var dailySales float64
		if err := s.db.Model(&internal.Order{}).
			Where("status IN ? AND created_at BETWEEN ? AND ?", validStatuses, startTimestamp, endTimestamp).
			Select("COALESCE(SUM(final_price), 0)").Scan(&dailySales).Error; err != nil {
			return nil, err
		}

		var orderCount int64
		if err := s.db.Model(&internal.Order{}).
			Where("created_at BETWEEN ? AND ?", startTimestamp, endTimestamp).
			Count(&orderCount).Error; err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"date":        date.Format("2006-01-02"),
			"sales":       dailySales,
			"order_count": orderCount,
		})
	}

	return results, nil
}

// GetTopProducts 获取热销商品统计
func (s *AdminStatsService) GetTopProducts(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 有效订单状态
	validStatuses := []string{
		internal.OrderStatusToSend,
		internal.OrderStatusToReceive,
		internal.OrderStatusFinished,
	}

	rows, err := s.db.Table("order_item").
		Select("sku.\"SPUID\" as spu_id, spu.name, spu.cover_image as image, SUM(order_item.quantity) as total_quantity, SUM(order_item.price * order_item.quantity) as total_sales").
		Joins("JOIN sku ON order_item.sku_id = sku.id").
		Joins("JOIN spu ON sku.\"SPUID\" = spu.id").
		Joins("JOIN \"order\" o ON order_item.order_id = o.id").
		Where("o.status IN ?", validStatuses).
		Group("sku.\"SPUID\", spu.name, spu.cover_image").
		Order("total_quantity DESC").
		Limit(limit).
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var spuID, name, image string
		var totalQuantity int64
		var totalSales float64

		if err := rows.Scan(&spuID, &name, &image, &totalQuantity, &totalSales); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"spu_id":         spuID,
			"name":           name,
			"image":          image,
			"total_quantity": totalQuantity,
			"total_sales":    totalSales,
		})
	}

	return results, nil
}

// GetUserRegistrationChart 获取用户注册图表数据
func (s *AdminStatsService) GetUserRegistrationChart(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.AddDate(0, 0, 1).Add(-time.Second)

		var userCount int64
		if err := s.db.Model(&internal.User{}).
			Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay).
			Count(&userCount).Error; err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"date":       date.Format("2006-01-02"),
			"user_count": userCount,
		})
	}

	return results, nil
}

// GetCategorySalesStats 获取分类销售统计
func (s *AdminStatsService) GetCategorySalesStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := s.db.Table("order_items").
		Select("c.id, c.name, SUM(order_items.quantity) as total_quantity, SUM(order_items.total_price) as total_sales").
		Joins("JOIN spus spu ON order_items.spu_id = spu.id").
		Joins("JOIN categories c ON spu.category_id = c.id").
		Joins("JOIN orders o ON order_items.order_id = o.id").
		Where("o.status IN ?", []string{"completed", "shipped"}).
		Group("c.id, c.name").
		Order("total_sales DESC").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var categoryID, name string
		var totalQuantity int64
		var totalSales float64

		if err := rows.Scan(&categoryID, &name, &totalQuantity, &totalSales); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"category_id":    categoryID,
			"name":           name,
			"total_quantity": totalQuantity,
			"total_sales":    totalSales,
		})
	}

	return results, nil
}

// GetOrderStatusStats 获取订单状态统计
func (s *AdminStatsService) GetOrderStatusStats() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	statuses := []string{
		"TO_PAY",
		"TO_SEND",
		"TO_RECEIVE",
		"FINISHED",
		"CANCELED",
		"RETURN_APPLIED",
		"RETURN_REFUSED",
		"RETURN_FINISH",
	}

	for _, status := range statuses {
		var count int64
		if err := s.db.Model(&internal.Order{}).Where("status = ?", status).Count(&count).Error; err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"status": status,
			"count":  count,
		})
	}

	return results, nil
}
