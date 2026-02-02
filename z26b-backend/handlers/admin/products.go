package admin

import (
	"net/http"
	"strconv"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// AdminGetProducts 获取商品列表
func (h *Handler) AdminGetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	categoryId := c.Query("categoryId")
	tagId := c.Query("tagId")
	status := c.Query("status")

	var products []internal.SPU
	var total int64

	query := h.DB.Model(&internal.SPU{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryId != "" {
		query = query.Where("category_id = ?", categoryId)
	}
	if tagId != "" {
		query = query.Where("id IN (SELECT spu_id FROM spu_tag WHERE tag_id = ?)", tagId)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products)

	// Load tags
	for i := range products {
		var spuTags []internal.SPUTag
		h.DB.Preload("Tag").Where("spu_id = ?", products[i].ID).Find(&spuTags)
		var tags []internal.Tag
		for _, st := range spuTags {
			if st.Tag != nil {
				tags = append(tags, *st.Tag)
			}
		}
		products[i].Tags = tags
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"list": products, "total": total, "page": page, "pageSize": pageSize},
	})
}

// AdminGetProduct 获取商品详情
func (h *Handler) AdminGetProduct(c *gin.Context) {
	id := c.Param("id")

	var product internal.SPU
	if err := h.DB.Preload("Category").First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var spuTags []internal.SPUTag
	h.DB.Preload("Tag").Where("spu_id = ?", id).Find(&spuTags)
	var tags []internal.Tag
	for _, st := range spuTags {
		if st.Tag != nil {
			tags = append(tags, *st.Tag)
		}
	}
	product.Tags = tags

	var skus []internal.SKU
	h.DB.Where(`"SPUID" = ?`, id).Find(&skus)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"product": product, "skus": skus}})
}

// AdminCreateProduct 创建商品
func (h *Handler) AdminCreateProduct(c *gin.Context) {
	var req struct {
		Name       string   `json:"name" binding:"required"`
		Detail     string   `json:"detail"`
		CoverImage string   `json:"coverImage"`
		Images     []string `json:"images"`
		CategoryID string   `json:"categoryId"`
		TagIDs     []string `json:"tagIds"`
		Tags       []struct {
			ID    string `json:"_id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"tags"`
		Status   string `json:"status"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写商品名称"})
		return
	}

	status := req.Status
	if status == "" {
		status = "ENABLED"
	}

	now := time.Now().UnixMilli()
	adminID := c.GetString("adminID")

	product := internal.SPU{
		ID: internal.GenerateUUID(), Name: req.Name, Detail: req.Detail,
		CoverImage: req.CoverImage, CategoryID: req.CategoryID,
		Status: status, Priority: req.Priority, CreatedAt: now, UpdatedAt: now,
		CreatedBy: adminID, UpdatedBy: adminID,
	}

	if req.Images != nil {
		product.SwipeImages = internal.ToJSON(req.Images)
	}

	tx := h.DB.Begin()
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// Create tags
	if len(req.Tags) > 0 {
		for _, tagData := range req.Tags {
			var tagId string
			if tagData.ID != "" {
				tagId = tagData.ID
				h.DB.Model(&internal.Tag{}).Where("id = ?", tagId).Updates(map[string]interface{}{
					"name": tagData.Name, "color": tagData.Color,
				})
			} else {
				tagId = internal.GenerateUUID()
				newTag := internal.Tag{
					ID: tagId, Name: tagData.Name, Color: tagData.Color,
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				tx.Create(&newTag)
			}
			spuTag := internal.SPUTag{
				ID: internal.GenerateUUID(), SPUID: product.ID, TagID: tagId, CreatedAt: time.Now(),
			}
			tx.Create(&spuTag)
		}
	} else if len(req.TagIDs) > 0 {
		for _, tagId := range req.TagIDs {
			spuTag := internal.SPUTag{
				ID: internal.GenerateUUID(), SPUID: product.ID, TagID: tagId, CreatedAt: time.Now(),
			}
			tx.Create(&spuTag)
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// AdminUpdateProduct 更新商品
func (h *Handler) AdminUpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product internal.SPU
	if err := h.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	var req struct {
		Name       string   `json:"name"`
		Detail     string   `json:"detail"`
		CoverImage string   `json:"coverImage"`
		Images     []string `json:"images"`
		CategoryID string   `json:"categoryId"`
		TagIDs     []string `json:"tagIds"`
		Tags       []struct {
			ID    string `json:"_id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"tags"`
		Status   string `json:"status"`
		Priority int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	adminID := c.GetString("adminID")
	updates := map[string]interface{}{"updated_at": time.Now().UnixMilli(), "updated_by": adminID}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Detail != "" {
		updates["detail"] = req.Detail
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.CategoryID != "" {
		updates["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Priority != 0 {
		updates["priority"] = req.Priority
	}
	if req.Images != nil {
		updates["swipe_images"] = internal.ToJSON(req.Images)
	}

	tx := h.DB.Begin()
	tx.Model(&product).Updates(updates)

	// Update tags - 总是先删除旧的关联
	tx.Where("spu_id = ?", id).Delete(&internal.SPUTag{})

	// 然后创建新的关联
	if len(req.Tags) > 0 {
		for _, tagData := range req.Tags {
			var tagId string
			if tagData.ID != "" {
				tagId = tagData.ID
				tx.Model(&internal.Tag{}).Where("id = ?", tagId).Updates(map[string]interface{}{
					"name": tagData.Name, "color": tagData.Color, "updated_at": time.Now(),
				})
			} else {
				tagId = internal.GenerateUUID()
				newTag := internal.Tag{
					ID: tagId, Name: tagData.Name, Color: tagData.Color,
					CreatedAt: time.Now(), UpdatedAt: time.Now(),
				}
				tx.Create(&newTag)
			}
			spuTag := internal.SPUTag{
				ID: internal.GenerateUUID(), SPUID: id, TagID: tagId, CreatedAt: time.Now(),
			}
			tx.Create(&spuTag)
		}
	} else if len(req.TagIDs) > 0 {
		for _, tagId := range req.TagIDs {
			spuTag := internal.SPUTag{
				ID: internal.GenerateUUID(), SPUID: id, TagID: tagId, CreatedAt: time.Now(),
			}
			tx.Create(&spuTag)
		}
	}

	tx.Commit()
	h.DB.Preload("Category").First(&product, "id = ?", id)

	var spuTags []internal.SPUTag
	h.DB.Preload("Tag").Where("spu_id = ?", id).Find(&spuTags)
	var tags []internal.Tag
	for _, st := range spuTags {
		if st.Tag != nil {
			tags = append(tags, *st.Tag)
		}
	}
	product.Tags = tags

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// AdminDeleteProduct 删除商品
func (h *Handler) AdminDeleteProduct(c *gin.Context) {
	id := c.Param("id")

	var product internal.SPU
	if err := h.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	tx := h.DB.Begin()
	tx.Where(`"SPUID" = ?`, id).Delete(&internal.SKU{})
	tx.Where("spu_id = ?", id).Delete(&internal.SPUTag{})
	tx.Delete(&product)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AdminToggleProductStatus 切换商品状态
func (h *Handler) AdminToggleProductStatus(c *gin.Context) {
	id := c.Param("id")

	var product internal.SPU
	if err := h.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		return
	}

	newStatus := "DISABLED"
	if product.Status == "DISABLED" {
		newStatus = "ENABLED"
	}

	h.DB.Model(&product).Updates(map[string]interface{}{
		"status": newStatus, "updatedAt": time.Now().UnixMilli(),
	})

	c.JSON(http.StatusOK, gin.H{"status": newStatus})
}

// AdminGetSKUs 获取商品SKU列表
func (h *Handler) AdminGetSKUs(c *gin.Context) {
	id := c.Param("id")
	var skus []internal.SKU
	h.DB.Where(`"SPUID" = ?`, id).Find(&skus)
	c.JSON(http.StatusOK, gin.H{"data": skus})
}

// AdminCreateSKU 创建SKU
func (h *Handler) AdminCreateSKU(c *gin.Context) {
	var req struct {
		SPUID       string  `json:"spuId" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Image       string  `json:"image"`
		Price       float64 `json:"price" binding:"required"`
		Count       int     `json:"count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写完整信息"})
		return
	}

	now := time.Now().UnixMilli()
	adminID := c.GetString("adminID")

	sku := internal.SKU{
		ID: internal.GenerateUUID(), SPUID: req.SPUID, Description: req.Description,
		Image: req.Image, Price: req.Price, Count: req.Count,
		CreatedAt: now, UpdatedAt: now, CreatedBy: adminID, UpdatedBy: adminID,
	}

	if err := h.DB.Create(&sku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建SKU失败: " + err.Error()})
		return
	}

	// 更新 SPU 价格范围
	h.updateSPUPriceRange(req.SPUID)

	c.JSON(http.StatusCreated, gin.H{"data": sku})
}

// AdminUpdateSKU 更新SKU
func (h *Handler) AdminUpdateSKU(c *gin.Context) {
	id := c.Param("id")

	var sku internal.SKU
	if err := h.DB.First(&sku, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SKU不存在"})
		return
	}

	var req struct {
		Description *string  `json:"description"`
		Image       *string  `json:"image"`
		Price       *float64 `json:"price"`
		Count       *int     `json:"count"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	adminID := c.GetString("adminID")
	now := time.Now().UnixMilli()
	updates := map[string]interface{}{}
	priceChanged := false

	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Image != nil {
		updates["image"] = *req.Image
	}
	if req.Price != nil {
		updates["price"] = *req.Price
		priceChanged = true
	}
	if req.Count != nil {
		updates["count"] = *req.Count
	}

	// 只有当有实际更新时才更新时间戳和操作人
	if len(updates) > 0 {
		updates["updated_at"] = now
		updates["updated_by"] = adminID
	}

	h.DB.Model(&sku).Updates(updates)
	h.DB.First(&sku, "id = ?", id)

	// 如果价格变更，更新 SPU 价格范围
	if priceChanged {
		h.updateSPUPriceRange(sku.SPUID)
	}

	c.JSON(http.StatusOK, gin.H{"data": sku})
}

// AdminDeleteSKU 删除SKU
func (h *Handler) AdminDeleteSKU(c *gin.Context) {
	id := c.Param("id")

	// 先获取 SKU 信息以便后续更新 SPU 价格
	var sku internal.SKU
	if err := h.DB.First(&sku, "id = ?", id).Error; err == nil {
		h.DB.Delete(&internal.SKU{}, "id = ?", id)
		// 更新 SPU 价格范围
		h.updateSPUPriceRange(sku.SPUID)
	} else {
		h.DB.Delete(&internal.SKU{}, "id = ?", id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// updateSPUPriceRange 更新 SPU 的价格范围
func (h *Handler) updateSPUPriceRange(spuID string) {
	h.DB.Exec(`
		UPDATE spu SET 
			min_price = COALESCE((SELECT MIN(price) FROM sku WHERE "SPUID" = ?), 0),
			max_price = COALESCE((SELECT MAX(price) FROM sku WHERE "SPUID" = ?), 0)
		WHERE id = ?
	`, spuID, spuID, spuID)
}
