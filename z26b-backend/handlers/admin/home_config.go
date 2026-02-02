package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============================================
// 首页配置 - 轮播图管理
// ============================================

// BannerItem 轮播图条目
type BannerItem struct {
	ID       string `json:"id"`
	Image    string `json:"image"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	LinkType string `json:"linkType"` // product, category, url, none
	Priority int    `json:"priority"`
}

// AdminGetBanners 获取轮播图列表
func (h *Handler) AdminGetBanners(c *gin.Context) {
	var swipers []internal.Swiper
	if err := h.DB.Order("priority ASC, created_at DESC").Find(&swipers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取轮播图失败"})
		return
	}

	// 转换为前端需要的格式
	banners := make([]BannerItem, 0, len(swipers))
	for _, s := range swipers {
		// 解析images JSON
		var images []string
		if s.Images != nil {
			json.Unmarshal(s.Images, &images)
		}

		image := ""
		if len(images) > 0 {
			image = images[0]
		}

		banners = append(banners, BannerItem{
			ID:       s.ID,
			Image:    image,
			Title:    s.Title,
			Link:     s.Link,
			LinkType: determineLinkType(s.Link),
			Priority: s.Priority,
		})
	}

	c.JSON(http.StatusOK, banners)
}

// AdminCreateBanner 创建轮播图
func (h *Handler) AdminCreateBanner(c *gin.Context) {
	var req BannerItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 生成ID
	id := uuid.New().String()

	// 图片转为JSON数组
	images := []string{}
	if req.Image != "" {
		images = append(images, req.Image)
	}
	imagesJSON, _ := json.Marshal(images)

	swiper := internal.Swiper{
		ID:        id,
		Images:    imagesJSON,
		Title:     req.Title,
		Link:      req.Link,
		Priority:  req.Priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.DB.Create(&swiper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建轮播图失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"_id":     id,
		"message": "创建成功",
	})
}

// AdminUpdateBanner 更新轮播图
func (h *Handler) AdminUpdateBanner(c *gin.Context) {
	id := c.Param("id")

	var swiper internal.Swiper
	if err := h.DB.First(&swiper, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "轮播图不存在"})
		return
	}

	var req BannerItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 图片转为JSON数组
	images := []string{}
	if req.Image != "" {
		images = append(images, req.Image)
	}
	imagesJSON, _ := json.Marshal(images)

	swiper.Images = imagesJSON
	swiper.Title = req.Title
	swiper.Link = req.Link
	swiper.Priority = req.Priority
	swiper.UpdatedAt = time.Now()

	if err := h.DB.Save(&swiper).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新轮播图失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// AdminDeleteBanner 删除轮播图
func (h *Handler) AdminDeleteBanner(c *gin.Context) {
	id := c.Param("id")

	if err := h.DB.Delete(&internal.Swiper{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AdminReorderBanners 重新排序轮播图
func (h *Handler) AdminReorderBanners(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 批量更新优先级
	for i, id := range req.IDs {
		h.DB.Model(&internal.Swiper{}).Where("id = ?", id).Update("priority", i)
	}

	c.JSON(http.StatusOK, gin.H{"message": "排序更新成功"})
}

// ============================================
// 首页配置 - 推荐商品管理（使用独立关联表）
// ============================================

// RecommendedProductResponse 推荐商品响应
type RecommendedProductResponse struct {
	ID           string   `json:"id"`
	ProductID    string   `json:"productId"`
	ProductName  string   `json:"productName"`
	ProductImage string   `json:"productImage"`
	ProductPrice float64  `json:"productPrice"`
	Tags         []string `json:"tags"`
	Priority     int      `json:"priority"`
}

// AdminGetRecommendedProducts 获取推荐商品列表
func (h *Handler) AdminGetRecommendedProducts(c *gin.Context) {
	var recommendations []internal.RecommendedProduct
	if err := h.DB.Preload("SPU").Order("priority ASC, created_at DESC").Find(&recommendations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取推荐商品失败"})
		return
	}

	products := make([]RecommendedProductResponse, 0, len(recommendations))
	for _, rec := range recommendations {
		// 解析标签
		var tags []string
		if rec.Tags != nil {
			json.Unmarshal(rec.Tags, &tags)
		}

		// 获取商品价格（从SKU）
		var minPrice float64
		h.DB.Model(&internal.SKU{}).Where("\"SPUID\" = ?", rec.SPUID).Select("COALESCE(MIN(price), 0)").Scan(&minPrice)

		productName := ""
		productImage := ""
		if rec.SPU != nil {
			productName = rec.SPU.Name
			productImage = rec.SPU.CoverImage
		}

		products = append(products, RecommendedProductResponse{
			ID:           rec.ID,
			ProductID:    rec.SPUID,
			ProductName:  productName,
			ProductImage: productImage,
			ProductPrice: minPrice,
			Tags:         tags,
			Priority:     rec.Priority,
		})
	}

	c.JSON(http.StatusOK, products)
}

// AdminAddRecommendedProduct 添加推荐商品
func (h *Handler) AdminAddRecommendedProduct(c *gin.Context) {
	var req struct {
		ProductID string   `json:"productId"`
		Tags      []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查商品是否存在
	var spu internal.SPU
	if err := h.DB.First(&spu, "id = ?", req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	// 检查是否已经推荐
	var existing internal.RecommendedProduct
	if err := h.DB.Where("spu_id = ?", req.ProductID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该商品已在推荐列表中"})
		return
	}

	// 获取当前最大优先级
	var maxPriority int
	h.DB.Model(&internal.RecommendedProduct{}).Select("COALESCE(MAX(priority), 0)").Scan(&maxPriority)

	// 标签转JSON
	tagsJSON, _ := json.Marshal(req.Tags)

	// 创建推荐记录
	rec := internal.RecommendedProduct{
		ID:        uuid.New().String(),
		SPUID:     req.ProductID,
		Tags:      tagsJSON,
		Priority:  maxPriority + 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.DB.Create(&rec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加推荐商品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加推荐商品成功", "_id": rec.ID})
}

// AdminUpdateRecommendedProduct 更新推荐商品标签
func (h *Handler) AdminUpdateRecommendedProduct(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Tags []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	tagsJSON, _ := json.Marshal(req.Tags)

	// 尝试通过推荐ID更新
	result := h.DB.Model(&internal.RecommendedProduct{}).Where("id = ?", id).Updates(map[string]interface{}{
		"tags":       tagsJSON,
		"updated_at": time.Now(),
	})

	// 如果没找到，尝试通过商品ID更新
	if result.RowsAffected == 0 {
		result = h.DB.Model(&internal.RecommendedProduct{}).Where("spu_id = ?", id).Updates(map[string]interface{}{
			"tags":       tagsJSON,
			"updated_at": time.Now(),
		})
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// AdminRemoveRecommendedProduct 移除推荐商品
func (h *Handler) AdminRemoveRecommendedProduct(c *gin.Context) {
	id := c.Param("id")

	// 尝试通过推荐ID删除
	result := h.DB.Delete(&internal.RecommendedProduct{}, "id = ?", id)

	// 如果没找到，尝试通过商品ID删除
	if result.RowsAffected == 0 {
		result = h.DB.Delete(&internal.RecommendedProduct{}, "spu_id = ?", id)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除成功"})
}

// AdminReorderRecommendedProducts 重新排序推荐商品
func (h *Handler) AdminReorderRecommendedProducts(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 批量更新优先级（支持推荐ID或商品ID）
	for i, id := range req.IDs {
		// 先尝试通过推荐ID更新
		result := h.DB.Model(&internal.RecommendedProduct{}).Where("id = ?", id).Update("priority", i)
		// 如果没找到，尝试通过商品ID更新
		if result.RowsAffected == 0 {
			h.DB.Model(&internal.RecommendedProduct{}).Where("spu_id = ?", id).Update("priority", i)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "排序更新成功"})
}

// ============================================
// 辅助函数
// ============================================

func determineLinkType(link string) string {
	if link == "" {
		return "none"
	}
	// 简单判断链接类型
	if len(link) > 0 && link[0] == '/' {
		// 内部链接
		if containsPath(link, "/pages/goods/") {
			return "product"
		}
		if containsPath(link, "/pages/category/") {
			return "category"
		}
	}
	return "url"
}

func containsPath(link, path string) bool {
	return len(link) >= len(path) && link[:len(path)] == path
}

// ============================================
// 商品搜索（用于选择推荐商品）
// ============================================

// AdminSearchProductsForRecommend 搜索商品（用于添加推荐）
func (h *Handler) AdminSearchProductsForRecommend(c *gin.Context) {
	keyword := c.Query("keyword")
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	// 获取已推荐的商品ID列表
	var recommendedIDs []string
	h.DB.Model(&internal.RecommendedProduct{}).Pluck("spu_id", &recommendedIDs)

	// 查询未推荐的商品
	query := h.DB.Model(&internal.SPU{})
	if len(recommendedIDs) > 0 {
		query = query.Where("id NOT IN ?", recommendedIDs)
	}

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var spus []internal.SPU
	if err := query.Limit(limit).Find(&spus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	// 简化返回数据，并获取价格
	results := make([]gin.H, 0, len(spus))
	for _, spu := range spus {
		// 获取最低价格
		var minPrice float64
		h.DB.Model(&internal.SKU{}).Where("\"SPUID\" = ?", spu.ID).Select("COALESCE(MIN(price), 0)").Scan(&minPrice)

		results = append(results, gin.H{
			"_id":        spu.ID,
			"name":       spu.Name,
			"coverImage": spu.CoverImage,
			"price":      minPrice,
		})
	}

	c.JSON(http.StatusOK, results)
}

// ============================================
// 首页内容（富文本）管理
// ============================================

// AdminGetHomeContents 获取所有首页内容
func (h *Handler) AdminGetHomeContents(c *gin.Context) {
	var contents []internal.HomeContent
	if err := h.DB.Order("priority ASC, created_at DESC").Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取首页内容失败"})
		return
	}
	c.JSON(http.StatusOK, contents)
}

// AdminGetHomeContent 获取单个首页内容
func (h *Handler) AdminGetHomeContent(c *gin.Context) {
	key := c.Param("key")

	var content internal.HomeContent
	if err := h.DB.Where("key = ?", key).First(&content).Error; err != nil {
		// 如果不存在，返回空内容
		c.JSON(http.StatusOK, gin.H{
			"_id":      "",
			"key":      key,
			"title":    "",
			"content":  "",
			"enabled":  true,
			"priority": 0,
		})
		return
	}
	c.JSON(http.StatusOK, content)
}

// AdminSaveHomeContent 保存首页内容（创建或更新）
func (h *Handler) AdminSaveHomeContent(c *gin.Context) {
	var req struct {
		Key      string `json:"key" binding:"required"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Enabled  *bool  `json:"enabled"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 查找是否已存在
	var existing internal.HomeContent
	err := h.DB.Where("key = ?", req.Key).First(&existing).Error

	if err != nil {
		// 不存在，创建新的
		content := internal.HomeContent{
			ID:        uuid.New().String(),
			Key:       req.Key,
			Title:     req.Title,
			Content:   req.Content,
			Enabled:   true,
			Priority:  req.Priority,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if req.Enabled != nil {
			content.Enabled = *req.Enabled
		}

		if err := h.DB.Create(&content).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
			return
		}
		c.JSON(http.StatusOK, content)
	} else {
		// 存在，更新
		existing.Title = req.Title
		existing.Content = req.Content
		existing.Priority = req.Priority
		existing.UpdatedAt = time.Now()
		if req.Enabled != nil {
			existing.Enabled = *req.Enabled
		}

		if err := h.DB.Save(&existing).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, existing)
	}
}

// AdminDeleteHomeContent 删除首页内容
func (h *Handler) AdminDeleteHomeContent(c *gin.Context) {
	key := c.Param("key")

	if err := h.DB.Where("key = ?", key).Delete(&internal.HomeContent{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
