package crm

import "z26b-backend/internal"

// CRMEventServiceInterface CRM事件服务接口
type CRMEventServiceInterface interface {
	// RecordEvent 记录CRM事件
	RecordEvent(event *internal.CRMEvent) error
	// GetEventsByUser 获取用户的事件列表
	GetEventsByUser(userID string, eventType string, page, pageSize int) ([]internal.CRMEvent, int64, error)
	// GetEventsBySPU 获取商品的事件列表
	GetEventsBySPU(spuID string, eventType string, page, pageSize int) ([]internal.CRMEvent, int64, error)
	// GetEventStats 获取事件统计
	GetEventStats(startTime, endTime int64) (map[string]int64, error)
	// GetDailyEventStats 获取每日事件统计
	GetDailyEventStats(days int) ([]map[string]interface{}, error)
}

// CustomerStatsServiceInterface 客户统计服务接口
type CustomerStatsServiceInterface interface {
	// GetOrCreateStats 获取或创建客户统计记录
	GetOrCreateStats(userID string) (*internal.CustomerStats, error)
	// UpdateStats 更新客户统计
	UpdateStats(stats *internal.CustomerStats) error
	// GetCustomerList 获取客户列表（带统计）
	GetCustomerList(page, pageSize int, sortBy string, sortOrder string) ([]internal.CustomerStats, int64, error)
	// GetTopCustomers 获取Top客户
	GetTopCustomers(limit int, orderBy string) ([]internal.CustomerStats, error)
	// GetCustomerLevelDistribution 获取客户等级分布
	GetCustomerLevelDistribution() (map[string]int64, error)
	// GetCustomerOverview 获取客户概览统计
	GetCustomerOverview() (map[string]interface{}, error)
	// RefreshCustomerStats 刷新客户统计
	RefreshCustomerStats(userID string) error
}

// ProductStatsServiceInterface 商品统计服务接口
type ProductStatsServiceInterface interface {
	// GetOrCreateStats 获取或创建商品统计记录
	GetOrCreateStats(spuID string) (*internal.ProductStats, error)
	// UpdateStats 更新商品统计
	UpdateStats(stats *internal.ProductStats) error
	// GetProductList 获取商品统计列表
	GetProductList(page, pageSize int, sortBy string, sortOrder string) ([]internal.ProductStats, int64, error)
	// GetTopProducts 获取Top商品
	GetTopProducts(limit int, orderBy string) ([]internal.ProductStats, error)
	// GetProductOverview 获取商品概览统计
	GetProductOverview() (map[string]interface{}, error)
	// GetCategoryStats 获取分类统计
	GetCategoryStats() ([]map[string]interface{}, error)
	// RefreshProductStats 刷新商品统计
	RefreshProductStats(spuID string) error
}
