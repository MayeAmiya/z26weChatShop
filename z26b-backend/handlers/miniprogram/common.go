package miniprogram

import (
	"z26b-backend/internal"
	"z26b-backend/services/crm"
	miniprogram_services "z26b-backend/services/miniprogram"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 小程序端处理器
type Handler struct {
	GoodsService    miniprogram_services.GoodsServiceInterface
	UserService     miniprogram_services.UserServiceInterface
	CartService     miniprogram_services.CartServiceInterface
	OrderService    miniprogram_services.OrderServiceInterface
	CommentService  miniprogram_services.CommentServiceInterface
	WechatService   miniprogram_services.WechatServiceInterface
	CRMEventService *crm.CRMEventService
	DB              *gorm.DB // 暂时保留，用于其他功能迁移
}

// NewHandler 创建处理器实例
func NewHandler(
	goodsService miniprogram_services.GoodsServiceInterface,
	userService miniprogram_services.UserServiceInterface,
	cartService miniprogram_services.CartServiceInterface,
	orderService miniprogram_services.OrderServiceInterface,
	commentService miniprogram_services.CommentServiceInterface,
	wechatService miniprogram_services.WechatServiceInterface,
	crmEventService *crm.CRMEventService,
	db *gorm.DB,
) *Handler {
	return &Handler{
		GoodsService:    goodsService,
		UserService:     userService,
		CartService:     cartService,
		OrderService:    orderService,
		CommentService:  commentService,
		WechatService:   wechatService,
		CRMEventService: crmEventService,
		DB:              db,
	}
}

// GetOrCreateUser 通过 OpenID 获取或创建用户
func (h *Handler) GetOrCreateUser(c *gin.Context) (*internal.User, error) {
	openID := c.GetString("openID")
	if openID == "" {
		openID = "oTest_dev_openid_001"
	}

	return h.UserService.GetOrCreateUser(openID)
}
