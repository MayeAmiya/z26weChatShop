package admin

import (
	"net/http"
	"strconv"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// AdminGetUsers 获取用户列表
func (h *Handler) AdminGetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	var users []internal.User
	var total int64

	query := h.DB.Model(&internal.User{})
	if keyword != "" {
		query = query.Where("nickName LIKE ? OR openid LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"list": users, "total": total, "page": page, "pageSize": pageSize},
	})
}

// AdminGetUser 获取用户详情
func (h *Handler) AdminGetUser(c *gin.Context) {
	id := c.Param("id")

	var user internal.User
	if err := h.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var orderCount int64
	h.DB.Model(&internal.Order{}).Where("userId = ?", id).Count(&orderCount)

	var totalSpent float64
	h.DB.Model(&internal.Order{}).Where("userId = ? AND status NOT IN ?", id, []string{internal.OrderStatusCanceled, internal.OrderStatusToPay}).
		Select("COALESCE(SUM(final_price), 0)").Scan(&totalSpent)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"user": user, "orderCount": orderCount, "totalSpent": totalSpent},
	})
}

// AdminGetUserOrders 获取用户订单
func (h *Handler) AdminGetUserOrders(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var orders []internal.Order
	var total int64

	query := h.DB.Model(&internal.Order{}).Where("userId = ?", id)
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Preload("Items.SKU").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"list": orders, "total": total, "page": page, "pageSize": pageSize},
	})
}

// AdminCreateTestUser 创建测试用户
func (h *Handler) AdminCreateTestUser(c *gin.Context) {
	var req struct {
		NickName string `json:"nickName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写用户昵称"})
		return
	}

	user := internal.User{
		ID:       internal.GenerateUUID(),
		OpenID:   "test_" + internal.GenerateUUID()[:8],
		NickName: req.NickName,
	}
	h.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"data": user})
}
