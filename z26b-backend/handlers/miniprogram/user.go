package miniprogram

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserInfo 获取用户信息
func (h *Handler) GetUserInfo(c *gin.Context) {
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUserInfo 更新用户信息
func (h *Handler) UpdateUserInfo(c *gin.Context) {
	user, err := h.GetOrCreateUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	var req struct {
		NickName string `json:"nickName"`
		Avatar   string `json:"avatar"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err = h.UserService.UpdateUserInfo(user.ID, req.NickName, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := h.UserService.GetUserInfo(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedUser})
}
