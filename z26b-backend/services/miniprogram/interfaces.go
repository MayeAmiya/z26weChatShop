package miniprogram

import (
	"z26b-backend/internal"
)

// GoodsService 商品服务接口
type GoodsServiceInterface interface {
	GetGoodsList(page, pageSize int, categoryID, search string) ([]internal.SPU, int64, error)
	GetGoodDetail(id string) (*internal.SPU, []internal.SKU, error)
	GetSKUDetail(id string) (*internal.SKU, error)
	GetSKUsBySpuID(spuID string) ([]internal.SKU, error)
	GetCategories() ([]internal.Category, error)
	SearchGoods(keyword string, page, pageSize int) ([]internal.SPU, int64, error)
	GetHomeSwiper() ([]internal.Swiper, error)
	GetHomeContent(key string) (*internal.HomeContent, error)
	GetHomeCategories() ([]internal.Category, error)
	GetPromotions() ([]internal.Promotion, error)
}

// UserService 用户服务接口
type UserServiceInterface interface {
	GetOrCreateUser(openID string) (*internal.User, error)
	GetUserByID(id string) (*internal.User, error)
	UpdateUserInfo(userID string, nickName, avatar string) error
	GetUserInfo(userID string) (*internal.User, error)
}

// CartService 购物车服务接口
type CartServiceInterface interface {
	GetCartItems(userID string) ([]internal.CartItem, error)
	AddToCart(userID, skuID string, quantity int) error
	UpdateCartItem(userID, itemID string, quantity int, isSelected *bool) error
	RemoveCartItem(userID, itemID string) error
	ClearCart(userID string) error
	GetCartItemCount(userID string) (int64, error)
}

// OrderService 订单服务接口
type OrderServiceInterface interface {
	GetOrderList(userID, status string, page, pageSize int) ([]internal.Order, int64, error)
	GetOrderDetail(orderID, userID string) (*internal.Order, error)
	CreateOrder(userID string, items []internal.OrderItem, addressID string) (*internal.Order, error)
	UpdateOrderStatus(orderID, userID, status string) error
	CancelOrder(orderID, userID string) error
	GetAdminOrderList(status, orderNo, userID string, page, pageSize int) ([]map[string]interface{}, int64, error)
	UpdateAdminOrderStatus(orderID, status string) error
}

// CommentService 评论服务接口
type CommentServiceInterface interface {
	GetGoodsComments(spuID string, page, pageSize int) ([]internal.Comment, int64, error)
	CreateComment(comment *internal.Comment) error
	GetCommentByID(id string) (*internal.Comment, error)
	UpdateComment(id string, updates map[string]interface{}) error
	DeleteComment(id string) error
	GetCommentsByUser(userID string, page, pageSize int) ([]internal.Comment, int64, error)
}

// AddressService 地址服务接口
type AddressServiceInterface interface {
	GetAddressList(userID string) ([]internal.Address, error)
	GetAddress(id, userID string) (*internal.Address, error)
	CreateAddress(address *internal.Address) error
	UpdateAddress(id, userID string, updates map[string]interface{}) error
	DeleteAddress(id, userID string) error
	SetDefaultAddress(id, userID string) error
}

// WechatService 微信服务接口
type WechatServiceInterface interface {
	WxLogin(code string) (map[string]interface{}, error)
	GenerateSign(params map[string]interface{}, key string) string
}
