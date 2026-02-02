package miniprogram

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WechatConfig 微信配置
type WechatConfig struct {
	AppID, AppSecret, MchID, MchKey, NotifyURL string
	Enabled                                    bool
}

// WxLogin 微信登录
func (h *Handler) WxLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := h.WechatService.WxLogin(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
