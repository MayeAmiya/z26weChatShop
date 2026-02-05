package crm

import (
	"time"

	"z26b-backend/internal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CustomerStatsService 客户统计服务
type CustomerStatsService struct {
	db *gorm.DB
}

// NewCustomerStatsService 创建客户统计服务实例
func NewCustomerStatsService(db *gorm.DB) *CustomerStatsService {
	return &CustomerStatsService{db: db}
}

// GetOrCreateStats 获取或创建客户统计记录
func (s *CustomerStatsService) GetOrCreateStats(userID string) (*internal.CustomerStats, error) {
	var stats internal.CustomerStats
	err := s.db.Where("user_id = ?", userID).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		stats = internal.CustomerStats{
			ID:            uuid.New().String(),
			UserID:        userID,
			CustomerLevel: "normal",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.db.Create(&stats).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpdateStats 更新客户统计
func (s *CustomerStatsService) UpdateStats(stats *internal.CustomerStats) error {
	stats.UpdatedAt = time.Now()
	return s.db.Save(stats).Error
}

// GetCustomerList 获取客户列表（带统计）
func (s *CustomerStatsService) GetCustomerList(page, pageSize int, sortBy string, sortOrder string) ([]internal.CustomerStats, int64, error) {
	var statsList []internal.CustomerStats
	var total int64

	query := s.db.Model(&internal.CustomerStats{}).Preload("User")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if sortBy == "" {
		sortBy = "total_spent"
	}
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	query = query.Order(sortBy + " " + sortOrder)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&statsList).Error; err != nil {
		return nil, 0, err
	}

	return statsList, total, nil
}

// GetTopCustomers 获取Top客户
func (s *CustomerStatsService) GetTopCustomers(limit int, orderBy string) ([]internal.CustomerStats, error) {
	var statsList []internal.CustomerStats

	if orderBy == "" {
		orderBy = "total_spent DESC"
	}

	err := s.db.Model(&internal.CustomerStats{}).
		Preload("User").
		Order(orderBy).
		Limit(limit).
		Find(&statsList).Error

	return statsList, err
}

// GetCustomerLevelDistribution 获取客户等级分布
func (s *CustomerStatsService) GetCustomerLevelDistribution() (map[string]int64, error) {
	distribution := make(map[string]int64)

	type Result struct {
		CustomerLevel string
		Count         int64
	}
	var results []Result

	err := s.db.Model(&internal.CustomerStats{}).
		Select("customer_level, COUNT(*) as count").
		Group("customer_level").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	for _, r := range results {
		distribution[r.CustomerLevel] = r.Count
	}

	return distribution, nil
}

// GetCustomerOverview 获取客户概览统计
func (s *CustomerStatsService) GetCustomerOverview() (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// 总客户数
	var totalCustomers int64
	s.db.Model(&internal.CustomerStats{}).Count(&totalCustomers)
	overview["totalCustomers"] = totalCustomers

	// 总消费金额
	var totalSpent float64
	s.db.Model(&internal.CustomerStats{}).Select("COALESCE(SUM(total_spent), 0)").Scan(&totalSpent)
	overview["totalSpent"] = totalSpent

	// 平均客单价
	var avgOrderValue float64
	s.db.Model(&internal.CustomerStats{}).Select("COALESCE(AVG(avg_order_value), 0)").Scan(&avgOrderValue)
	overview["avgOrderValue"] = avgOrderValue

	// 活跃客户数（30天内有活动）
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).UnixMilli()
	var activeCustomers int64
	s.db.Model(&internal.CustomerStats{}).Where("last_active_at >= ?", thirtyDaysAgo).Count(&activeCustomers)
	overview["activeCustomers"] = activeCustomers

	return overview, nil
}

// RefreshCustomerStats 刷新客户统计（根据事件重新计算）
func (s *CustomerStatsService) RefreshCustomerStats(userID string) error {
	stats, err := s.GetOrCreateStats(userID)
	if err != nil {
		return err
	}

	// 从订单表统计
	var orderStats struct {
		TotalOrders int
		TotalSpent  float64
	}
	s.db.Model(&internal.Order{}).
		Where("user_id = ? AND status NOT IN ?", userID, []string{internal.OrderStatusCanceled, internal.OrderStatusReturnFinish}).
		Select("COUNT(*) as total_orders, COALESCE(SUM(final_price), 0) as total_spent").
		Scan(&orderStats)

	stats.TotalOrders = orderStats.TotalOrders
	stats.TotalSpent = orderStats.TotalSpent
	if orderStats.TotalOrders > 0 {
		stats.AvgOrderValue = orderStats.TotalSpent / float64(orderStats.TotalOrders)
	}

	// 统计退款
	var refundStats struct {
		TotalRefunds int
		RefundAmount float64
	}
	s.db.Model(&internal.Order{}).
		Where("user_id = ? AND status = ?", userID, internal.OrderStatusReturnFinish).
		Select("COUNT(*) as total_refunds, COALESCE(SUM(final_price), 0) as refund_amount").
		Scan(&refundStats)
	stats.TotalRefunds = refundStats.TotalRefunds
	stats.RefundAmount = refundStats.RefundAmount

	// 从CRM事件表统计
	var viewCount int64
	s.db.Model(&internal.CRMEvent{}).Where("user_id = ? AND event_type = ?", userID, internal.CRMEventTypeView).Count(&viewCount)
	stats.TotalViews = int(viewCount)

	var cartCount int64
	s.db.Model(&internal.CRMEvent{}).Where("user_id = ? AND event_type = ?", userID, internal.CRMEventTypeCart).Count(&cartCount)
	stats.TotalCarts = int(cartCount)

	var commentCount int64
	s.db.Model(&internal.Comment{}).Where("user_id = ?", userID).Count(&commentCount)
	stats.TotalComments = int(commentCount)

	// 最后下单时间
	var lastOrder internal.Order
	if s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&lastOrder).Error == nil {
		stats.LastOrderAt = &lastOrder.CreatedAt
	}

	// 最后活跃时间
	var lastEvent internal.CRMEvent
	if s.db.Where("user_id = ?", userID).Order("created_at DESC").First(&lastEvent).Error == nil {
		stats.LastActiveAt = &lastEvent.CreatedAt
	}

	// 更新客户等级
	stats.CustomerLevel = calculateCustomerLevel(stats.TotalSpent, stats.TotalOrders)

	return s.UpdateStats(stats)
}

// calculateCustomerLevel 根据消费情况计算客户等级
func calculateCustomerLevel(totalSpent float64, totalOrders int) string {
	if totalSpent >= 10000 || totalOrders >= 50 {
		return "diamond"
	} else if totalSpent >= 5000 || totalOrders >= 30 {
		return "platinum"
	} else if totalSpent >= 2000 || totalOrders >= 15 {
		return "gold"
	} else if totalSpent >= 500 || totalOrders >= 5 {
		return "silver"
	}
	return "normal"
}
