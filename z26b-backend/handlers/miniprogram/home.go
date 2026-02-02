package miniprogram

import (
	"net/http"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// GetHomeSwiper 获取首页轮播图
func (h *Handler) GetHomeSwiper(c *gin.Context) {
	swipers, err := h.GoodsService.GetHomeSwiper()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch swipers"})
		return
	}
	if len(swipers) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": []internal.Swiper{}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": swipers[0]})
}

// GetHomeContent 获取首页富文本内容
func (h *Handler) GetHomeContent(c *gin.Context) {
	key := c.DefaultQuery("key", "main")
	content, err := h.GoodsService.GetHomeContent(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": content})
}

// GetHomeCategories 获取首页分类
func (h *Handler) GetHomeCategories(c *gin.Context) {
	categories, err := h.GoodsService.GetHomeCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetPromotions 获取促销活动
func (h *Handler) GetPromotions(c *gin.Context) {
	promotions, err := h.GoodsService.GetPromotions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch promotions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": promotions})
}
