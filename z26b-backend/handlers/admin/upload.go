package admin

import (
	"net/http"
	"path/filepath"
	"strings"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

var allowedImageTypes = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
var allowedMimeTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true}

// AdminUploadImage 上传单张图片
func (h *Handler) AdminUploadImage(c *gin.Context) {
	if !internal.IsMinIOInitialized() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "存储服务未初始化"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageTypes[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的图片格式"})
		return
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的图片类型"})
		return
	}

	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过 5MB"})
		return
	}

	url, err := internal.UploadFile(c.Request.Context(), file, header.Filename, contentType, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url, "filename": header.Filename, "size": header.Size})
}

// AdminUploadImages 上传多张图片
func (h *Handler) AdminUploadImages(c *gin.Context) {
	if !internal.IsMinIOInitialized() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "存储服务未初始化"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "一次最多上传 10 张图片"})
		return
	}

	var results []gin.H
	var errors []string

	for _, header := range files {
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if !allowedImageTypes[ext] {
			errors = append(errors, header.Filename+": 不支持的格式")
			continue
		}

		contentType := header.Header.Get("Content-Type")
		if !allowedMimeTypes[contentType] {
			errors = append(errors, header.Filename+": 无效的类型")
			continue
		}

		if header.Size > 5*1024*1024 {
			errors = append(errors, header.Filename+": 超过 5MB")
			continue
		}

		file, err := header.Open()
		if err != nil {
			errors = append(errors, header.Filename+": 打开失败")
			continue
		}

		url, err := internal.UploadFile(c.Request.Context(), file, header.Filename, contentType, header.Size)
		file.Close()

		if err != nil {
			errors = append(errors, header.Filename+": 上传失败")
			continue
		}

		results = append(results, gin.H{"url": url, "filename": header.Filename, "size": header.Size})
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "errors": errors})
}

// AdminDeleteImage 删除图片
func (h *Handler) AdminDeleteImage(c *gin.Context) {
	if !internal.IsMinIOInitialized() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "存储服务未初始化"})
		return
	}

	var req struct {
		URL string `json:"url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供图片URL"})
		return
	}

	if err := internal.DeleteFile(c.Request.Context(), req.URL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
