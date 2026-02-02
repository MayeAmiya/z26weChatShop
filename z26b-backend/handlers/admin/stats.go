package admin

import (
	"net/http"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// AdminGetStats 获取统计数据
func (h *Handler) AdminGetStats(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).UnixMilli()

	var todayOrders, totalProducts, totalUsers, pendingOrders, toShipOrders int64
	var todaySales, monthlySales float64

	h.DB.Model(&internal.Order{}).Where("created_at >= ?", todayStart).Count(&todayOrders)
	h.DB.Model(&internal.Order{}).Where("created_at >= ? AND status != ?", todayStart, internal.OrderStatusToPay).
		Select("COALESCE(SUM(final_price), 0)").Scan(&todaySales)
	h.DB.Model(&internal.SPU{}).Count(&totalProducts)
	h.DB.Model(&internal.User{}).Count(&totalUsers)
	h.DB.Model(&internal.Order{}).Where("status = ?", internal.OrderStatusToPay).Count(&pendingOrders)
	h.DB.Model(&internal.Order{}).Where("status = ?", internal.OrderStatusToSend).Count(&toShipOrders)

	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).UnixMilli()
	h.DB.Model(&internal.Order{}).Where("created_at >= ? AND status NOT IN ?", monthStart, []string{internal.OrderStatusToPay, internal.OrderStatusCanceled}).
		Select("COALESCE(SUM(final_price), 0)").Scan(&monthlySales)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"todayOrders": todayOrders, "todaySales": todaySales, "totalProducts": totalProducts,
			"totalUsers": totalUsers, "pendingOrders": pendingOrders, "toShipOrders": toShipOrders,
			"monthlySales": monthlySales,
		},
	})
}

// AdminGetSalesStats 获取销售统计
func (h *Handler) AdminGetSalesStats(c *gin.Context) {
	days := 7
	now := time.Now()

	type DailyStat struct {
		Date   string  `json:"date"`
		Orders int64   `json:"orders"`
		Sales  float64 `json:"sales"`
	}

	var stats []DailyStat
	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).UnixMilli()
		dayEnd := dayStart + 86400000

		var orders int64
		var sales float64
		h.DB.Model(&internal.Order{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&orders)
		h.DB.Model(&internal.Order{}).Where("created_at >= ? AND created_at < ? AND status NOT IN ?", dayStart, dayEnd, []string{internal.OrderStatusToPay, internal.OrderStatusCanceled}).
			Select("COALESCE(SUM(final_price), 0)").Scan(&sales)

		stats = append(stats, DailyStat{Date: date.Format("01-02"), Orders: orders, Sales: sales})
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// AdminGetTopProducts 获取热销商品
func (h *Handler) AdminGetTopProducts(c *gin.Context) {
	type TopProduct struct {
		SPUID     string  `json:"spuId"`
		Name      string  `json:"name"`
		SoldCount int     `json:"soldCount"`
		Sales     float64 `json:"sales"`
	}

	var results []struct {
		SKUID     string
		SoldCount int
		Sales     float64
	}

	h.DB.Model(&internal.OrderItem{}).
		Select("sku_id as skuid, SUM(quantity) as sold_count, SUM(price * quantity) as sales").
		Group("sku_id").Order("sold_count DESC").Limit(10).Scan(&results)

	var topProducts []TopProduct
	for _, r := range results {
		var sku internal.SKU
		if h.DB.Preload("SPU").First(&sku, "id = ?", r.SKUID).Error == nil && sku.SPU != nil {
			topProducts = append(topProducts, TopProduct{
				SPUID: sku.SPUID, Name: sku.SPU.Name, SoldCount: r.SoldCount, Sales: r.Sales,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": topProducts})
}

// AdminGetOrderStatusStats 获取订单状态统计
func (h *Handler) AdminGetOrderStatusStats(c *gin.Context) {
	type StatusStat struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	statuses := []string{internal.OrderStatusToPay, internal.OrderStatusToSend, internal.OrderStatusToReceive, internal.OrderStatusFinished, internal.OrderStatusCanceled}
	var stats []StatusStat
	for _, status := range statuses {
		var count int64
		h.DB.Model(&internal.Order{}).Where("status = ?", status).Count(&count)
		stats = append(stats, StatusStat{Status: status, Count: count})
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
