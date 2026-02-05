package miniprogram

import (
	"net/http"
	"strconv"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// GetGoodsList 获取商品列表
func (h *Handler) GetGoodsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	categoryID := c.Query("categoryId")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	goods, total, err := h.GoodsService.GetGoodsList(page, pageSize, categoryID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch goods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"records": goods, "total": total}})
}

// GetGood 获取商品详情
func (h *Handler) GetGood(c *gin.Context) {
	id := c.Param("id")

	good, skus, err := h.GoodsService.GetGoodDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
		return
	}

	// 记录商品浏览事件
	go func() {
		user, _ := h.GetOrCreateUser(c)
		userID := ""
		if user != nil {
			userID = user.ID
		}
		h.CRMEventService.RecordEvent(&internal.CRMEvent{
			UserID:    userID,
			EventType: internal.CRMEventTypeView,
			SPUID:     id,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}()

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"_id": good.ID, "name": good.Name, "detail": good.Detail,
			"cover_image": good.CoverImage, "swiper_images": good.SwipeImages,
			"status": good.Status, "skus": skus,
		},
	})
}

// GetSKU 获取SKU详情
func (h *Handler) GetSKU(c *gin.Context) {
	id := c.Param("id")

	sku, err := h.GoodsService.GetSKUDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sku})
}

// GetSKUsBySpuId 获取商品的所有SKU
func (h *Handler) GetSKUsBySpuId(c *gin.Context) {
	spuID := c.Param("spuId")

	skus, err := h.GoodsService.GetSKUsBySpuID(spuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SKUs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": skus})
}

// GetCategories 获取分类列表
func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.GoodsService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// SearchGoods 搜索商品
func (h *Handler) SearchGoods(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	goods, total, err := h.GoodsService.SearchGoods(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search goods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"records": goods, "total": total}})
}
