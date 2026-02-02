package admin

import (
	"net/http"
	"strconv"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminGetOrders 获取订单列表
func (h *Handler) AdminGetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	orderNo := c.Query("orderNo")
	userId := c.Query("userId")

	var orders []internal.Order
	var total int64

	query := h.DB.Model(&internal.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if orderNo != "" {
		query = query.Where("id LIKE ?", "%"+orderNo+"%")
	}
	if userId != "" {
		query = query.Where("userId = ?", userId)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Preload("Items.SKU.SPU").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	type OrderWithUser struct {
		internal.Order
		User *internal.User `json:"user"`
	}

	var result []OrderWithUser
	for _, order := range orders {
		owu := OrderWithUser{Order: order}
		var user internal.User
		if h.DB.First(&user, "id = ?", order.UserID).Error == nil {
			owu.User = &user
		}
		result = append(result, owu)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"list": result, "total": total, "page": page, "pageSize": pageSize},
	})
}

// AdminGetOrder 获取订单详情
func (h *Handler) AdminGetOrder(c *gin.Context) {
	id := c.Param("id")

	var order internal.Order
	if err := h.DB.Preload("Items.SKU.SPU").First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var user internal.User
	h.DB.First(&user, "id = ?", order.UserID)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"order": order, "user": user}})
}

// AdminShipOrder 发货
func (h *Handler) AdminShipOrder(c *gin.Context) {
	id := c.Param("id")

	var order internal.Order
	if err := h.DB.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if order.Status != internal.OrderStatusToSend {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不正确，无法发货"})
		return
	}

	order.Status = internal.OrderStatusToReceive
	order.UpdatedAt = time.Now().UnixMilli()
	h.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "发货成功"})
}

// AdminRefundOrder 退款
func (h *Handler) AdminRefundOrder(c *gin.Context) {
	id := c.Param("id")

	var order internal.Order
	if err := h.DB.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if order.Status != internal.OrderStatusToSend && order.Status != internal.OrderStatusToReceive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态不正确，无法退款"})
		return
	}

	tx := h.DB.Begin()

	// 恢复库存
	var items []internal.OrderItem
	h.DB.Where("order_id = ?", id).Find(&items)
	for _, item := range items {
		tx.Model(&internal.SKU{}).Where("id = ?", item.SKUID).Update("count", gorm.Expr("count + ?", item.Quantity))
	}

	order.Status = internal.OrderStatusReturnFinish
	order.UpdatedAt = time.Now().UnixMilli()
	tx.Save(&order)

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "退款成功", "refundAmount": order.FinalPrice})
}

// AdminUpdateOrderStatus 更新订单状态
func (h *Handler) AdminUpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择状态"})
		return
	}

	var order internal.Order
	if err := h.DB.First(&order, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	order.Status = req.Status
	order.UpdatedAt = time.Now().UnixMilli()
	h.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}
