package miniprogram

import (
	"net/http"
	"strconv"
	"time"

	"z26b-backend/internal"

	"github.com/gin-gonic/gin"
)

// GetGoodsComments 获取商品评论
func (h *Handler) GetGoodsComments(c *gin.Context) {
	spuID := c.Param("spuId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	comments, total, err := h.CommentService.GetGoodsComments(spuID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{"pageNum": page, "pageSize": pageSize, "totalCount": total, "pageList": comments},
	})
}

// SubmitComment 提交评论
func (h *Handler) SubmitComment(c *gin.Context) {
	var req struct {
		SPUID          string `json:"spuId" binding:"required"`
		SKUID          string `json:"skuId"`
		UserName       string `json:"userName"`
		CommentContent string `json:"commentContent" binding:"required"`
		CommentScore   int    `json:"commentScore" binding:"required,min=1,max=5"`
		IsAnonymity    bool   `json:"isAnonymity"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.DefaultQuery("userId", "ANONYMOUS")

	comment := &internal.Comment{
		ID:             internal.GenerateUUID(),
		SPUID:          req.SPUID,
		SKUID:          req.SKUID,
		UserID:         userID,
		UserName:       req.UserName,
		CommentContent: req.CommentContent,
		CommentScore:   req.CommentScore,
		IsAnonymity:    req.IsAnonymity,
		CreatedAt:      time.Now().Unix(),
	}

	err := h.CommentService.CreateComment(comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comment})
}
