package miniprogram

import (
	"net/http"
	"strconv"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// GetOrderList 获取订单列表
func (h *Handler) GetOrderList(c *gin.Context) {
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")

	orders, total, err := h.OrderService.GetOrderList(user.ID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"records": orders, "total": total, "page": page, "pageSize": pageSize},
	})
}

// GetOrderDetail 获取订单详情
func (h *Handler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	order, err := h.OrderService.GetOrderDetail(id, user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// CreateOrder 创建订单
func (h *Handler) CreateOrder(c *gin.Context) {
	var req struct {
		AddressID string `json:"addressId" binding:"required"`
		Remarks   string `json:"remarks"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// 获取选中的购物车商品
	cartItems, err := h.CartService.GetCartItems(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart items"})
		return
	}

	// 过滤选中的商品
	var selectedItems []internal.CartItem
	for _, item := range cartItems {
		if item.IsSelected {
			selectedItems = append(selectedItems, item)
		}
	}

	if len(selectedItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No items selected"})
		return
	}

	// 转换为订单项
	var orderItems []internal.OrderItem
	for _, cartItem := range selectedItems {
		if cartItem.SKU != nil {
			orderItem := internal.OrderItem{
				ID:       internal.GenerateUUID(),
				SKUID:    cartItem.SKUID,
				Quantity: cartItem.Quantity,
				Price:    cartItem.SKU.Price,
			}
			orderItems = append(orderItems, orderItem)
		}
	}

	order, err := h.OrderService.CreateOrder(user.ID, orderItems, req.AddressID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录购买事件
	go func() {
		for _, item := range selectedItems {
			if item.SKU != nil {
				h.CRMEventService.RecordEvent(&internal.CRMEvent{
					UserID:    user.ID,
					EventType: internal.CRMEventTypePurchase,
					SPUID:     item.SKU.SPUID,
					SKUID:     item.SKUID,
					OrderID:   order.ID,
					Amount:    item.SKU.Price * float64(item.Quantity),
					IPAddress: c.ClientIP(),
					UserAgent: c.Request.UserAgent(),
				})
			}
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"order":         order,
			"paidAmount":    order.FinalPrice,
			"paymentMethod": "DIRECT",
		},
	})
}

// CancelOrder 取消订单
func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = h.OrderService.CancelOrder(id, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to cancel order"})
		return
	}

	order, err := h.OrderService.GetOrderDetail(id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order, "message": "Order canceled successfully"})
}

// ConfirmReceipt 确认收货
func (h *Handler) ConfirmReceipt(c *gin.Context) {
	id := c.Param("id")
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = h.OrderService.UpdateOrderStatus(id, user.ID, internal.OrderStatusFinished)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to confirm receipt"})
		return
	}

	order, err := h.OrderService.GetOrderDetail(id, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
