package crm

import (
	"time"

	"z26b-backend/internal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductStatsService 商品统计服务
type ProductStatsService struct {
	db *gorm.DB
}

// NewProductStatsService 创建商品统计服务实例
func NewProductStatsService(db *gorm.DB) *ProductStatsService {
	return &ProductStatsService{db: db}
}

// GetOrCreateStats 获取或创建商品统计记录
func (s *ProductStatsService) GetOrCreateStats(spuID string) (*internal.ProductStats, error) {
	var stats internal.ProductStats
	err := s.db.Where("spu_id = ?", spuID).First(&stats).Error
	if err == gorm.ErrRecordNotFound {
		stats = internal.ProductStats{
			ID:        uuid.New().String(),
			SPUID:     spuID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.db.Create(&stats).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &stats, nil
}

// UpdateStats 更新商品统计
func (s *ProductStatsService) UpdateStats(stats *internal.ProductStats) error {
	stats.UpdatedAt = time.Now()
	return s.db.Save(stats).Error
}

// GetProductList 获取商品统计列表
func (s *ProductStatsService) GetProductList(page, pageSize int, sortBy string, sortOrder string) ([]internal.ProductStats, int64, error) {
	var statsList []internal.ProductStats
	var total int64

	query := s.db.Model(&internal.ProductStats{}).Preload("SPU")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if sortBy == "" {
		sortBy = "total_revenue"
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

// GetTopProducts 获取Top商品
func (s *ProductStatsService) GetTopProducts(limit int, orderBy string) ([]internal.ProductStats, error) {
	var statsList []internal.ProductStats

	if orderBy == "" {
		orderBy = "total_revenue DESC"
	}

	err := s.db.Model(&internal.ProductStats{}).
		Preload("SPU").
		Order(orderBy).
		Limit(limit).
		Find(&statsList).Error

	return statsList, err
}

// GetProductOverview 获取商品概览统计
func (s *ProductStatsService) GetProductOverview() (map[string]interface{}, error) {
	overview := make(map[string]interface{})

	// 总商品数
	var totalProducts int64
	s.db.Model(&internal.ProductStats{}).Count(&totalProducts)
	overview["totalProducts"] = totalProducts

	// 总浏览量
	var totalViews int64
	s.db.Model(&internal.ProductStats{}).Select("COALESCE(SUM(total_views), 0)").Scan(&totalViews)
	overview["totalViews"] = totalViews

	// 总销量
	var totalSales int64
	s.db.Model(&internal.ProductStats{}).Select("COALESCE(SUM(total_sales), 0)").Scan(&totalSales)
	overview["totalSales"] = totalSales

	// 总营收
	var totalRevenue float64
	s.db.Model(&internal.ProductStats{}).Select("COALESCE(SUM(total_revenue), 0)").Scan(&totalRevenue)
	overview["totalRevenue"] = totalRevenue

	// 平均转化率
	var avgConversionRate float64
	s.db.Model(&internal.ProductStats{}).Select("COALESCE(AVG(conversion_rate), 0)").Scan(&avgConversionRate)
	overview["avgConversionRate"] = avgConversionRate

	// 平均评分
	var avgScore float64
	s.db.Model(&internal.ProductStats{}).Select("COALESCE(AVG(avg_score), 0)").Scan(&avgScore)
	overview["avgScore"] = avgScore

	return overview, nil
}

// GetCategoryStats 获取分类统计
func (s *ProductStatsService) GetCategoryStats() ([]map[string]interface{}, error) {
	type CategoryStat struct {
		CategoryID   string
		CategoryName string
		TotalSales   int64
		TotalRevenue float64
	}
	var results []CategoryStat

	err := s.db.Model(&internal.ProductStats{}).
		Select("spu.category_id, category.name as category_name, SUM(product_stats.total_sales) as total_sales, SUM(product_stats.total_revenue) as total_revenue").
		Joins("LEFT JOIN spu ON product_stats.spu_id = spu.id").
		Joins("LEFT JOIN category ON spu.category_id = category.id").
		Group("spu.category_id, category.name").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	var stats []map[string]interface{}
	for _, r := range results {
		stats = append(stats, map[string]interface{}{
			"categoryId":   r.CategoryID,
			"categoryName": r.CategoryName,
			"totalSales":   r.TotalSales,
			"totalRevenue": r.TotalRevenue,
		})
	}

	return stats, nil
}

// RefreshProductStats 刷新商品统计（根据事件重新计算）
func (s *ProductStatsService) RefreshProductStats(spuID string) error {
	stats, err := s.GetOrCreateStats(spuID)
	if err != nil {
		return err
	}

	// 从CRM事件表统计浏览量
	var viewCount int64
	s.db.Model(&internal.CRMEvent{}).Where("spu_id = ? AND event_type = ?", spuID, internal.CRMEventTypeView).Count(&viewCount)
	stats.TotalViews = int(viewCount)

	// 从CRM事件表统计加购数
	var cartCount int64
	s.db.Model(&internal.CRMEvent{}).Where("spu_id = ? AND event_type = ?", spuID, internal.CRMEventTypeCart).Count(&cartCount)
	stats.TotalCarts = int(cartCount)

	// 从订单项统计销量和营收
	var salesStats struct {
		TotalSales   int
		TotalRevenue float64
	}
	s.db.Model(&internal.OrderItem{}).
		Joins("JOIN sku ON order_item.sku_id = sku.id").
		Joins("JOIN \"order\" ON order_item.order_id = \"order\".id").
		Where("sku.\"SPUID\" = ? AND \"order\".status NOT IN ?", spuID, []string{internal.OrderStatusCanceled, internal.OrderStatusReturnFinish}).
		Select("COALESCE(SUM(order_item.quantity), 0) as total_sales, COALESCE(SUM(order_item.price * order_item.quantity), 0) as total_revenue").
		Scan(&salesStats)
	stats.TotalSales = salesStats.TotalSales
	stats.TotalRevenue = salesStats.TotalRevenue

	// 统计退款
	var refundStats struct {
		TotalRefunds int
		RefundAmount float64
	}
	s.db.Model(&internal.OrderItem{}).
		Joins("JOIN sku ON order_item.sku_id = sku.id").
		Joins("JOIN \"order\" ON order_item.order_id = \"order\".id").
		Where("sku.\"SPUID\" = ? AND \"order\".status = ?", spuID, internal.OrderStatusReturnFinish).
		Select("COALESCE(SUM(order_item.quantity), 0) as total_refunds, COALESCE(SUM(order_item.price * order_item.quantity), 0) as refund_amount").
		Scan(&refundStats)
	stats.TotalRefunds = refundStats.TotalRefunds
	stats.RefundAmount = refundStats.RefundAmount

	// 统计评论数和平均评分
	var commentStats struct {
		TotalComments int64
		AvgScore      float64
	}
	s.db.Model(&internal.Comment{}).
		Where("spu_id = ?", spuID).
		Select("COUNT(*) as total_comments, COALESCE(AVG(comment_score), 0) as avg_score").
		Scan(&commentStats)
	stats.TotalComments = int(commentStats.TotalComments)
	stats.AvgScore = commentStats.AvgScore

	// 计算转化率 (销量 / 浏览量)
	if stats.TotalViews > 0 {
		stats.ConversionRate = float64(stats.TotalSales) / float64(stats.TotalViews)
	} else {
		stats.ConversionRate = 0
	}

	return s.UpdateStats(stats)
}
