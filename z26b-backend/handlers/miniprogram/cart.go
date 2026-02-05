package miniprogram

import (
	"net/http"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// GetCartItems 获取购物车列表
func (h *Handler) GetCartItems(c *gin.Context) {
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	items, err := h.CartService.GetCartItems(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"isNotEmpty": len(items) > 0,
			"storeGoods": []map[string]interface{}{
				{"storeId": "1000", "storeName": "Default Store", "storeStatus": 1, "goodsList": items},
			},
		},
	})
}

// AddToCart 添加到购物车
func (h *Handler) AddToCart(c *gin.Context) {
	var req struct {
		SKUID    string `json:"skuId" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
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

	err = h.CartService.AddToCart(user.ID, req.SKUID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 记录加购事件
	go func() {
		// 获取SKU关联的SPU
		var sku internal.SKU
		if h.DB.First(&sku, "id = ?", req.SKUID).Error == nil {
			h.CRMEventService.RecordEvent(&internal.CRMEvent{
				UserID:    user.ID,
				EventType: internal.CRMEventTypeCart,
				SPUID:     sku.SPUID,
				SKUID:     req.SKUID,
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
			})
		}
	}()

	// 获取更新后的购物车项
	items, err := h.CartService.GetCartItems(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated cart"})
		return
	}

	// 找到刚添加的项
	var addedItem *internal.CartItem
	for _, item := range items {
		if item.SKUID == req.SKUID {
			addedItem = &item
			break
		}
	}

	if addedItem != nil {
		c.JSON(http.StatusOK, gin.H{"data": addedItem})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find added item"})
	}
}

// UpdateCartItem 更新购物车项
func (h *Handler) UpdateCartItem(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Quantity   *int  `json:"quantity,omitempty"`
		IsSelected *bool `json:"isSelected,omitempty"`
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

	// 如果提供了数量，使用它，否则使用 -1 表示不更新数量
	quantity := -1
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	err = h.CartService.UpdateCartItem(user.ID, id, quantity, req.IsSelected)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// 获取更新后的项
	items, err := h.CartService.GetCartItems(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated cart"})
		return
	}

	// 找到更新的项
	var updatedItem *internal.CartItem
	for _, item := range items {
		if item.ID == id {
			updatedItem = &item
			break
		}
	}

	if updatedItem != nil {
		c.JSON(http.StatusOK, gin.H{"data": updatedItem})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
	}
}

// RemoveFromCart 从购物车移除
func (h *Handler) RemoveFromCart(c *gin.Context) {
	id := c.Param("id")

	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = h.CartService.RemoveCartItem(user.ID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed"})
}

// ClearCart 清空购物车
func (h *Handler) ClearCart(c *gin.Context) {
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	err = h.CartService.ClearCart(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
}
