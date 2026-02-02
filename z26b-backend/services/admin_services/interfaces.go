package admin_services

import (
	"time"

	"z26b-backend/internal"
)

// AdminGoodsService 管理后台商品服务接口
type AdminGoodsServiceInterface interface {
	GetProducts(page, pageSize int, keyword, categoryID, tagID, status string) ([]internal.SPU, int64, error)
	GetProduct(id string) (*internal.SPU, []internal.SKU, error)
	CreateProduct(spu *internal.SPU, skus []internal.SKU) error
	UpdateProduct(id string, updates map[string]interface{}) error
	DeleteProduct(id string) error
	UpdateProductStatus(id, status string) error
}

// AdminCategoryService 管理后台分类服务接口
type AdminCategoryServiceInterface interface {
	GetCategories() ([]internal.Category, error)
	GetCategory(id string) (*internal.Category, error)
	CreateCategory(req *internal.Category) error
	UpdateCategory(id string, updates map[string]interface{}) error
	DeleteCategory(id string) error
}

// AdminUserManagementService 管理后台管理员用户管理服务接口
type AdminUserManagementServiceInterface interface {
	GetAdminUsers(page, pageSize int, keyword, status string) ([]internal.Admin, int64, error)
	GetAdminUserByID(adminUserID string) (*internal.Admin, error)
	CreateAdminUser(username, password, role string) (*internal.Admin, error)
	UpdateAdminUser(adminUserID string, updates map[string]interface{}) error
	DeleteAdminUser(adminUserID string) error
	AuthenticateAdmin(username, password string) (*internal.Admin, error)
}

// AdminOrderService 管理后台订单服务接口
type AdminOrderServiceInterface interface {
	GetOrders(page, pageSize int, keyword, status, userID string, startDate, endDate *time.Time) ([]internal.Order, int64, error)
	GetOrderByID(orderID string) (*internal.Order, error)
	UpdateOrderStatus(orderID, status string) error
	UpdateOrderShipping(orderID, shippingCompany, trackingNumber string) error
	GetOrderStats() (map[string]interface{}, error)
	RefundOrder(orderID string, refundAmount float64, refundReason string) error
	ConfirmRefund(orderID string) error
}

// AdminTagService 管理后台标签服务接口
type AdminTagServiceInterface interface {
	GetTags(page, pageSize int, keyword, status string) ([]internal.Tag, int64, error)
	GetTagByID(tagID string) (*internal.Tag, error)
	CreateTag(name, description, color string, sortOrder int) (*internal.Tag, error)
	UpdateTag(tagID string, updates map[string]interface{}) error
	DeleteTag(tagID string) error
	UpdateTagStatus(tagID, status string) error
	GetAllActiveTags() ([]internal.Tag, error)
}

// AdminStatsService 管理后台统计服务接口
type AdminStatsServiceInterface interface {
	GetStats() (map[string]interface{}, error)
	GetSalesStats(days int) ([]map[string]interface{}, error)
	GetTopProducts(limit int) ([]map[string]interface{}, error)
	GetOrderStatusStats() ([]map[string]interface{}, error)
}

// AdminAuthService 管理后台认证服务接口
type AdminAuthServiceInterface interface {
	Login(username, password string) (*internal.Admin, error)
	Logout(adminUserID string) error
	ChangePassword(adminUserID, oldPassword, newPassword string) error
	ResetPassword(adminUserID, newPassword string) error
	ValidateToken(token string) (*internal.Admin, error)
	GetCurrentUser(adminUserID string) (*internal.Admin, error)
	CheckPermission(adminUserID, permission string) (bool, error)
	GetUserPermissions(adminUserID string) ([]string, error)
}