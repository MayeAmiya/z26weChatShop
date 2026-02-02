package admin

import (
	"net/http"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// AdminGetTags 获取标签列表
func (h *Handler) AdminGetTags(c *gin.Context) {
	var tags []internal.Tag
	query := h.DB.Model(&internal.Tag{})

	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query.Order("sort DESC, created_at DESC").Find(&tags)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tags})
}

// AdminGetTag 获取标签详情
func (h *Handler) AdminGetTag(c *gin.Context) {
	id := c.Param("id")
	var tag internal.Tag
	if err := h.DB.First(&tag, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "标签不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tag})
}

// AdminCreateTag 创建标签
func (h *Handler) AdminCreateTag(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		SortOrder   int    `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": err.Error()})
		return
	}

	var count int64
	h.DB.Model(&internal.Tag{}).Where("name = ?", input.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "标签名称已存在"})
		return
	}

	tag := internal.Tag{
		ID:          internal.GenerateUUID(),
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
		SortOrder:   input.SortOrder,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if tag.Color == "" {
		tag.Color = "#1890ff"
	}

	h.DB.Create(&tag)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tag})
}

// AdminUpdateTag 更新标签
func (h *Handler) AdminUpdateTag(c *gin.Context) {
	id := c.Param("id")

	var tag internal.Tag
	if err := h.DB.First(&tag, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "error": "标签不存在"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		SortOrder   int    `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": err.Error()})
		return
	}

	if input.Name != "" && input.Name != tag.Name {
		var count int64
		h.DB.Model(&internal.Tag{}).Where("name = ? AND id != ?", input.Name, id).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "标签名称已存在"})
			return
		}
		tag.Name = input.Name
	}

	if input.Color != "" {
		tag.Color = input.Color
	}
	if input.Description != "" {
		tag.Description = input.Description
	}
	if input.SortOrder > 0 {
		tag.SortOrder = input.SortOrder
	}
	tag.UpdatedAt = time.Now()

	h.DB.Save(&tag)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": tag})
}

// AdminDeleteTag 删除标签
func (h *Handler) AdminDeleteTag(c *gin.Context) {
	id := c.Param("id")

	h.DB.Where("tag_id = ?", id).Delete(&internal.SPUTag{})
	h.DB.Delete(&internal.Tag{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "删除成功"})
}
