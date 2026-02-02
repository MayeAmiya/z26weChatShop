package admin

import (
	"net/http"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// AdminGetCategories 获取分类列表
func (h *Handler) AdminGetCategories(c *gin.Context) {
	categories, err := h.AdminCategoryService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// AdminGetCategory 获取分类详情
func (h *Handler) AdminGetCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.AdminCategoryService.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": category})
}

// AdminCreateCategory 创建分类
func (h *Handler) AdminCreateCategory(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Icon     string `json:"icon"`
		Image    string `json:"image"`
		ParentID string `json:"parentId"`
		Sort     int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写分类名称"})
		return
	}

	category := &internal.Category{
		Name:     req.Name,
		Icon:     req.Icon,
		Image:    req.Image,
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}

	err := h.AdminCategoryService.CreateCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

// AdminUpdateCategory 更新分类
func (h *Handler) AdminUpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name     string `json:"name"`
		Icon     string `json:"icon"`
		Image    string `json:"image"`
		ParentID string `json:"parentId"`
		Sort     int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Image != "" {
		updates["image"] = req.Image
	}
	if req.ParentID != "" {
		updates["parentId"] = req.ParentID
	}
	updates["sort"] = req.Sort

	err := h.AdminCategoryService.UpdateCategory(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	// 获取更新后的分类
	category, err := h.AdminCategoryService.GetCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// AdminDeleteCategory 删除分类
func (h *Handler) AdminDeleteCategory(c *gin.Context) {
	id := c.Param("id")

	err := h.AdminCategoryService.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类下有商品，无法删除"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
