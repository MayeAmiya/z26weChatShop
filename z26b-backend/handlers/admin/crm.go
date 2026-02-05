package admin

import (
	"net/http"
	"strconv"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// ============================================
// 商品分析 API
// ============================================

// AdminGetProductAnalysisOverview 获取商品分析概览
func (h *Handler) AdminGetProductAnalysisOverview(c *gin.Context) {
	overview, err := h.ProductStatsService.GetProductOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品概览失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// AdminGetProductStatsList 获取商品统计列表
func (h *Handler) AdminGetProductStatsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "total_revenue")
	sortOrder := c.DefaultQuery("sortOrder", "DESC")

	list, total, err := h.ProductStatsService.GetProductList(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品统计列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":     list,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminGetTopProducts 获取热销/热门商品
func (h *Handler) AdminGetTopProductStats(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("orderBy", "total_revenue DESC")

	products, err := h.ProductStatsService.GetTopProducts(limit, orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取Top商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// AdminRefreshProductStats 刷新商品统计
func (h *Handler) AdminRefreshProductStats(c *gin.Context) {
	spuID := c.Param("id")
	if spuID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "商品ID不能为空"})
		return
	}

	if err := h.ProductStatsService.RefreshProductStats(spuID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刷新商品统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "刷新成功"})
}

// AdminRefreshAllProductStats 刷新所有商品统计
func (h *Handler) AdminRefreshAllProductStats(c *gin.Context) {
	// 获取所有商品
	var spus []internal.SPU
	if err := h.DB.Find(&spus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品列表失败"})
		return
	}

	successCount := 0
	for _, spu := range spus {
		if err := h.ProductStatsService.RefreshProductStats(spu.ID); err == nil {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "刷新完成",
		"total":   len(spus),
		"success": successCount,
	})
}

// ============================================
// 客户分析 API
// ============================================

// AdminGetCustomerAnalysisOverview 获取客户分析概览
func (h *Handler) AdminGetCustomerAnalysisOverview(c *gin.Context) {
	overview, err := h.CustomerStatsService.GetCustomerOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取客户概览失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": overview})
}

// AdminGetCustomerStatsList 获取客户统计列表
func (h *Handler) AdminGetCustomerStatsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	sortBy := c.DefaultQuery("sortBy", "total_spent")
	sortOrder := c.DefaultQuery("sortOrder", "DESC")

	list, total, err := h.CustomerStatsService.GetCustomerList(page, pageSize, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取客户统计列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":     list,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminGetTopCustomers 获取高价值客户
func (h *Handler) AdminGetTopCustomers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orderBy := c.DefaultQuery("orderBy", "total_spent DESC")

	customers, err := h.CustomerStatsService.GetTopCustomers(limit, orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取高价值客户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// AdminGetCustomerLevelDistribution 获取客户等级分布
func (h *Handler) AdminGetCustomerLevelDistribution(c *gin.Context) {
	distribution, err := h.CustomerStatsService.GetCustomerLevelDistribution()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取客户等级分布失败"})
		return
	}

	// 转换为前端需要的格式
	var result []gin.H
	totalCount := int64(0)
	for _, count := range distribution {
		totalCount += count
	}

	levels := []string{"normal", "silver", "gold", "platinum", "diamond", "vip"}
	for _, level := range levels {
		count := distribution[level]
		percentage := float64(0)
		if totalCount > 0 {
			percentage = float64(count) / float64(totalCount) * 100
		}
		result = append(result, gin.H{
			"level":      level,
			"count":      count,
			"percentage": percentage,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// AdminRefreshCustomerStats 刷新客户统计
func (h *Handler) AdminRefreshCustomerStats(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	if err := h.CustomerStatsService.RefreshCustomerStats(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刷新客户统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "刷新成功"})
}

// AdminRefreshAllCustomerStats 刷新所有客户统计
func (h *Handler) AdminRefreshAllCustomerStats(c *gin.Context) {
	// 获取所有用户
	var users []internal.User
	if err := h.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	successCount := 0
	for _, user := range users {
		if err := h.CustomerStatsService.RefreshCustomerStats(user.ID); err == nil {
			successCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "刷新完成",
		"total":   len(users),
		"success": successCount,
	})
}

// ============================================
// CRM 事件 API
// ============================================

// AdminGetCRMEventStats 获取CRM事件统计
func (h *Handler) AdminGetCRMEventStats(c *gin.Context) {
	startTimeStr := c.Query("startTime")
	endTimeStr := c.Query("endTime")

	var startTime, endTime int64
	if startTimeStr != "" {
		startTime, _ = strconv.ParseInt(startTimeStr, 10, 64)
	}
	if endTimeStr != "" {
		endTime, _ = strconv.ParseInt(endTimeStr, 10, 64)
	}

	stats, err := h.CRMEventService.GetEventStats(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取事件统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// AdminGetCRMEventsByUser 获取用户的事件列表
func (h *Handler) AdminGetCRMEventsByUser(c *gin.Context) {
	userID := c.Param("userId")
	eventType := c.Query("eventType")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	events, total, err := h.CRMEventService.GetEventsByUser(userID, eventType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户事件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":     events,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminGetCRMEventsBySPU 获取商品的事件列表
func (h *Handler) AdminGetCRMEventsBySPU(c *gin.Context) {
	spuID := c.Param("spuId")
	eventType := c.Query("eventType")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	events, total, err := h.CRMEventService.GetEventsBySPU(spuID, eventType, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品事件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":     events,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// AdminGetDailyEventStats 获取每日事件统计
func (h *Handler) AdminGetDailyEventStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	stats, err := h.CRMEventService.GetDailyEventStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取每日事件统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// AdminGetCRMDashboard 获取CRM仪表盘数据（汇总）
func (h *Handler) AdminGetCRMDashboard(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()

	// 今日事件统计
	todayStats, _ := h.CRMEventService.GetEventStats(todayStart, 0)

	// 商品概览
	productOverview, _ := h.ProductStatsService.GetProductOverview()

	// 客户概览
	customerOverview, _ := h.CustomerStatsService.GetCustomerOverview()

	// Top 5 商品
	topProducts, _ := h.ProductStatsService.GetTopProducts(5, "total_revenue DESC")

	// Top 5 客户
	topCustomers, _ := h.CustomerStatsService.GetTopCustomers(5, "total_spent DESC")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"todayEvents":      todayStats,
			"productOverview":  productOverview,
			"customerOverview": customerOverview,
			"topProducts":      topProducts,
			"topCustomers":     topCustomers,
		},
	})
}
